package workers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	dbTypes "github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy/pkg/commands/artifact"
	"github.com/aquasecurity/trivy/pkg/commands/auth"
	"github.com/aquasecurity/trivy/pkg/db"
	"github.com/aquasecurity/trivy/pkg/fanal/analyzer"
	ftypes "github.com/aquasecurity/trivy/pkg/fanal/types"
	"github.com/aquasecurity/trivy/pkg/flag"
	"github.com/aquasecurity/trivy/pkg/report"
	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/helpers"
	models "github.com/estebangarcia/cm3070-final-project/pkg/oci_models"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/misconfiguration"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/vulnerability"
	"github.com/google/go-containerregistry/pkg/name"
	"golang.org/x/sync/errgroup"
)

type SecurityScannerWorker struct {
	cfg                *config.AppConfig
	parallelScans      int
	manifestRepository *repositories.ManifestRepository
	s3Client           *s3.Client
}

func NewSecurityScannerWorker(parallelScans int, s3Client *s3.Client, manifestRepository *repositories.ManifestRepository, cfg *config.AppConfig) *SecurityScannerWorker {
	return &SecurityScannerWorker{
		cfg:                cfg,
		parallelScans:      parallelScans,
		manifestRepository: manifestRepository,
		s3Client:           s3Client,
	}
}

func (w *SecurityScannerWorker) Handle(ctx context.Context) error {
	err := auth.Login(ctx, w.cfg.BaseURL, flag.Options{
		RegistryOptions: flag.RegistryOptions{
			Credentials: []ftypes.Credential{
				{
					Username: w.cfg.AdminUser.Email,
					Password: w.cfg.AdminUser.Password,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	eg, grpCtx := errgroup.WithContext(ctx)
	eg.SetLimit(w.parallelScans)

	manifests, err := w.manifestRepository.GetAllUnscanned(ctx)
	if err != nil {
		return err
	}

	for _, manifest := range manifests {
		eg.Go(w.startSecurityScan(grpCtx, manifest))
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (w *SecurityScannerWorker) startSecurityScan(ctx context.Context, manifest *ent.Manifest) func() error {
	return func() error {
		opts, err := getBaseOpts()
		if err != nil {
			return err
		}

		r, err := artifact.NewRunner(ctx, *opts)
		if err != nil {
			return err
		}

		defer r.Close(ctx)

		scans := map[string]func(context.Context, flag.Options) (types.Report, error){
			"application/vnd.docker.container.image.v1+json": r.ScanImage,
			"application/vnd.oci.image.config.v1+json":       r.ScanImage,
			"application/vnd.python.artifact":                r.ScanFilesystem,
			"application/vnd.cncf.helm.config.v1+json":       r.ScanFilesystem,
		}

		prepareFuncMapping := map[string]func(context.Context, *ent.Manifest, *flag.Options) error{
			"application/vnd.docker.container.image.v1+json": w.prepareContainerScan,
			"application/vnd.oci.image.config.v1+json":       w.prepareContainerScan,
			"application/vnd.python.artifact":                w.preparePythonScan,
			"application/vnd.cncf.helm.config.v1+json":       w.prepareHelmScan,
		}

		cleanupFuncMapping := map[string]func(*flag.Options) error{
			"application/vnd.docker.container.image.v1+json": w.cleanupContainerScan,
			"application/vnd.oci.image.config.v1+json":       w.cleanupContainerScan,
			"application/vnd.python.artifact":                w.cleanupPythonScan,
			"application/vnd.cncf.helm.config.v1+json":       w.cleanupPythonScan,
		}

		if prepareFunc, ok := prepareFuncMapping[manifest.ArtifactType]; ok && prepareFunc(ctx, manifest, opts) != nil {
			return err
		}

		defer func() {
			cleanupFunc, ok := cleanupFuncMapping[manifest.ArtifactType]
			if ok {
				fmt.Printf("Running cleanup function for %s", manifest.ArtifactType)
				cleanupFunc(opts)
			}
		}()

		scan, ok := scans[manifest.ArtifactType]
		if !ok {
			fmt.Printf("Manifest type %s for %s not supported. Skipping\n", manifest.ArtifactType, opts.ScanOptions.Target)
			if err := w.manifestRepository.MarkAsScanned(ctx, manifest); err != nil {
				return err
			}
			return nil
		}

		fmt.Printf("Starting scan for %s\n", opts.ScanOptions.Target)
		rawReport, err := scan(ctx, *opts)
		if err != nil {
			return err
		}

		vulnReport, err := r.Filter(ctx, *opts, rawReport)
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)

		writer := &report.JSONWriter{
			Output:         buf,
			ListAllPkgs:    true,
			ShowSuppressed: false,
		}

		if err = writer.Write(ctx, vulnReport); err != nil {
			return err
		}

		var report models.TrivyReport

		if err = json.Unmarshal(buf.Bytes(), &report); err != nil {
			return err
		}

		var vulnerabilities ent.Vulnerabilities

		var misconfigurations ent.Misconfigurations
		var misconfigIDs []string

		for _, results := range report.Results {
			for _, rawVulnerability := range results.Vulnerabilities {
				vulnerabilities = append(vulnerabilities, &ent.Vulnerability{
					PackageName:             rawVulnerability.PackageName,
					VulnerabilityID:         rawVulnerability.VulnerabilityID,
					VulnerabilityURLDetails: rawVulnerability.PrimaryURL,
					InstalledVersion:        rawVulnerability.InstalledVersion,
					FixedVersion:            rawVulnerability.FixedVersion,
					Status:                  vulnerability.Status(rawVulnerability.Status),
					Severity:                vulnerability.Severity(rawVulnerability.Severity),
					Title:                   rawVulnerability.Title,
				})
			}

			for _, rawMisconfiguration := range results.Misconfigurations {
				misconfigurations = append(misconfigurations, &ent.Misconfiguration{
					MisconfigurationID:         rawMisconfiguration.ID,
					MisconfigurationURLDetails: rawMisconfiguration.PrimaryURL,
					Title:                      rawMisconfiguration.Title,
					Severity:                   misconfiguration.Severity(rawMisconfiguration.Severity),
				})
				if !slices.Contains(misconfigIDs, rawMisconfiguration.ID) {
					misconfigIDs = append(misconfigIDs, rawMisconfiguration.ID)
				}
			}
		}

		if err = w.manifestRepository.CreateMisconfigurationsInBulk(ctx, misconfigurations); err != nil {
			return err
		}

		misconfigs, err := w.manifestRepository.GetMisconfigurationsByIDs(ctx, misconfigIDs)
		if err != nil {
			return err
		}

		var manifestMisconfigurations ent.ManifestMisconfigurations
		targetFileMisconfigs := map[string][]string{}

		misconfigMap := make(map[string]*ent.Misconfiguration)
		for _, m := range misconfigs {
			misconfigMap[m.MisconfigurationID] = m
		}

		for _, results := range report.Results {
			for _, rawMisconfiguration := range results.Misconfigurations {
				misconfig, ok := misconfigMap[rawMisconfiguration.ID]
				if !ok {
					continue
				}

				tFileMisconfigs, ok := targetFileMisconfigs[results.Target]
				if !ok {
					tFileMisconfigs = []string{}
				}

				if !slices.Contains(tFileMisconfigs, misconfig.MisconfigurationID) {
					manifestMisconfigurations = append(manifestMisconfigurations, &ent.ManifestMisconfiguration{
						TargetFile:         results.Target,
						Message:            rawMisconfiguration.Message,
						Resolution:         rawMisconfiguration.Resolution,
						MisconfigurationID: misconfig.ID,
						ManifestID:         manifest.ID,
					})
					tFileMisconfigs = append(tFileMisconfigs, misconfig.MisconfigurationID)
					targetFileMisconfigs[results.Target] = tFileMisconfigs
				}
			}
		}

		if err = w.manifestRepository.CreateManifestMisconfigurationsInBulk(ctx, manifestMisconfigurations); err != nil {
			return err
		}

		if err = w.manifestRepository.CreateVulnerabilitiesInBulkAndMarkAsScanned(ctx, vulnerabilities, manifest); err != nil {
			return err
		}

		fmt.Printf("Scan for %s finished\n", opts.ScanOptions.Target)

		return nil
	}
}

func (w *SecurityScannerWorker) prepareContainerScan(ctx context.Context, manifest *ent.Manifest, flags *flag.Options) error {
	organizaxtionSlug := manifest.Edges.Repository.Edges.Registry.Edges.Organization.Slug
	registrySlug := manifest.Edges.Repository.Edges.Registry.Slug
	repositoryName := manifest.Edges.Repository.Name
	digest := manifest.Digest
	flags.ScanOptions.Target = fmt.Sprintf("%s/%s/%s/%s@%s", w.cfg.BaseURL, organizaxtionSlug, registrySlug, repositoryName, digest)

	return nil
}

func (w *SecurityScannerWorker) cleanupContainerScan(flags *flag.Options) error {
	return nil
}

func (w *SecurityScannerWorker) cleanupPythonScan(flags *flag.Options) error {
	return os.Remove(flags.ScanOptions.Target)
}

func (w *SecurityScannerWorker) prepareHelmScan(ctx context.Context, manifest *ent.Manifest, flags *flag.Options) error {
	var chartLayer *ent.ManifestLayer

	for _, layer := range manifest.Edges.ManifestLayers {
		if layer.MediaType == "application/vnd.cncf.helm.chart.content.v1.tar+gzip" {
			chartLayer = layer
			break
		}
	}

	if chartLayer == nil {
		return fmt.Errorf("helm chart manifest with digest '%s' has no helm chart content layer", manifest.Digest)
	}

	output, err := w.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &w.cfg.S3.BlobsBucketName,
		Key:    aws.String(helpers.GetS3KeyForBlob(manifest.Edges.Repository.Edges.Registry.Edges.Organization.Slug, chartLayer.Digest)),
	})

	dirName, err := os.MkdirTemp("/tmp", "helm-chart-scanning*")
	if err != nil {
		return err
	}

	gzipReader, err := gzip.NewReader(output.Body)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzipReader)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		destinationPath := filepath.Join(dirName, header.Name)
		dir := filepath.Dir(destinationPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) && os.MkdirAll(dir, 0755) != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			outFile, err := os.Create(destinationPath)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			outFile.Close()
		}
	}

	flags.ScanOptions.Target = dirName + "/" + manifest.Edges.Repository.Name
	flags.MisconfOptions.HelmKubeVersion = "1.32"
	flags.ScanOptions.Scanners = types.Scanners{
		types.MisconfigScanner,
	}
	flags.DisabledAnalyzers = append(analyzer.TypeOSes, analyzer.TypeLanguages...)
	flags.PackageOptions = flag.PackageOptions{}

	return nil
}

func (w *SecurityScannerWorker) preparePythonScan(ctx context.Context, manifest *ent.Manifest, flags *flag.Options) error {
	var whlLayer *ent.ManifestLayer

	for _, layer := range manifest.Edges.ManifestLayers {
		if layer.MediaType == "application/vnd.python.whl.file" {
			whlLayer = layer
			break
		}
	}

	if whlLayer == nil {
		return fmt.Errorf("python manifest with digest '%s' has no whl layer", manifest.Digest)
	}

	requirementsFileContents := fmt.Sprintf("%s==%s\n", whlLayer.Annotations["Name"], whlLayer.Annotations["Version"])

	if requiresDist, ok := whlLayer.Annotations["Requires-Dist"]; ok {
		var whlDeps []string
		err := json.Unmarshal([]byte(requiresDist), &whlDeps)
		if err != nil {
			return err
		}

		for _, dep := range whlDeps {
			requirementsFileContents = requirementsFileContents + dep + "\n"
		}
	}

	dirName, err := os.MkdirTemp("/tmp", "python-scanning*")
	if err != nil {
		return err
	}

	tempFileName := dirName + "/requirements.txt"

	file, err := os.Create(tempFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(requirementsFileContents); err != nil {
		return err
	}

	flags.ScanOptions.Target = tempFileName
	return nil
}

func getBaseOpts() (*flag.Options, error) {
	var dbRepositories []name.Reference
	for _, repo := range []string{db.DefaultGCRRepository, db.DefaultGHCRRepository} {
		ref, err := parseRepository(repo, db.SchemaVersion)
		if err != nil {
			return nil, err
		}
		dbRepositories = append(dbRepositories, ref)
	}

	return &flag.Options{
		GlobalOptions: flag.GlobalOptions{
			Quiet: true,
		},
		DBOptions: flag.DBOptions{
			DBRepositories: dbRepositories,
		},
		MisconfOptions: flag.MisconfOptions{
			MisconfigScanners: analyzer.TypeConfigFiles,
		},
		ScanOptions: flag.ScanOptions{
			Scanners: types.Scanners{
				types.VulnerabilityScanner,
				types.SecretScanner,
			},
		},
		ImageOptions: flag.ImageOptions{
			ImageSources: ftypes.AllImageSources,
		},
		PackageOptions: flag.PackageOptions{
			PkgTypes:         types.PkgTypes,
			PkgRelationships: ftypes.Relationships,
		},
		ReportOptions: flag.ReportOptions{
			Severities: toSeverity(dbTypes.SeverityNames),
		},
	}, nil
}

func parseRepository(repo string, dbSchemaVersion int) (name.Reference, error) {
	dbRepository, err := name.ParseReference(repo, name.WithDefaultTag(""))
	if err != nil {
		return nil, err
	}

	// Add the schema version if the tag is not specified for backward compatibility.
	t, ok := dbRepository.(name.Tag)
	if !ok || t.TagStr() != "" {
		return dbRepository, nil
	}

	dbRepository = t.Tag(strconv.Itoa(dbSchemaVersion))
	return dbRepository, nil
}

func toSeverity(severity []string) []dbTypes.Severity {
	if len(severity) == 0 {
		return nil
	}

	var severities []dbTypes.Severity

	for _, severityName := range severity {
		s, err := dbTypes.NewSeverity(severityName)
		if err == nil {
			severities = append(severities, s)
		}
	}

	return severities
}
