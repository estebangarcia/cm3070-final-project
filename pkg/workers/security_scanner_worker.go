package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	dbTypes "github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy/pkg/commands/artifact"
	"github.com/aquasecurity/trivy/pkg/commands/auth"
	"github.com/aquasecurity/trivy/pkg/db"
	ftypes "github.com/aquasecurity/trivy/pkg/fanal/types"
	"github.com/aquasecurity/trivy/pkg/flag"
	"github.com/aquasecurity/trivy/pkg/report"
	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	models "github.com/estebangarcia/cm3070-final-project/pkg/oci_models"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/vulnerability"
	"github.com/google/go-containerregistry/pkg/name"
	"golang.org/x/sync/errgroup"
)

type SecurityScannerWorker struct {
	cfg                *config.AppConfig
	parallelScans      int
	manifestRepository *repositories.ManifestRepository
}

func NewSecurityScannerWorker(parallelScans int, manifestRepository *repositories.ManifestRepository, cfg *config.AppConfig) *SecurityScannerWorker {
	return &SecurityScannerWorker{
		cfg:                cfg,
		parallelScans:      parallelScans,
		manifestRepository: manifestRepository,
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
		eg.Go(func() error {
			opts, err := getOpts(w.cfg.BaseURL, manifest.Edges.Repository.Edges.Registry.Edges.Organization.Slug, manifest.Edges.Repository.Edges.Registry.Slug, manifest.Edges.Repository.Name, manifest.Digest)
			if err != nil {
				return err
			}

			r, err := artifact.NewRunner(grpCtx, *opts)
			if err != nil {
				return err
			}

			defer r.Close(grpCtx)

			scans := map[string]func(context.Context, flag.Options) (types.Report, error){
				"application/vnd.docker.container.image.v1+json": r.ScanImage,
				"application/vnd.oci.image.config.v1+json":       r.ScanImage,
				//"application/vnd.python.artifact":                r.ScanFilesystem,
			}

			scan, ok := scans[manifest.ArtifactType]
			if !ok {
				fmt.Printf("Manifest type %s for %s not supported. Skipping\n", manifest.ArtifactType, opts.ScanOptions.Target)
				if err := w.manifestRepository.MarkAsScanned(grpCtx, manifest); err != nil {
					return err
				}
				return nil
			}

			fmt.Printf("Starting scan for %s\n", opts.ScanOptions.Target)
			rawReport, err := scan(grpCtx, *opts)
			if err != nil {
				return err
			}

			vulnReport, err := r.Filter(grpCtx, *opts, rawReport)
			if err != nil {
				return err
			}

			buf := new(bytes.Buffer)

			writer := &report.JSONWriter{
				Output:         buf,
				ListAllPkgs:    true,
				ShowSuppressed: false,
			}

			if err = writer.Write(grpCtx, vulnReport); err != nil {
				return err
			}

			var report models.TrivyReport

			if err = json.Unmarshal(buf.Bytes(), &report); err != nil {
				return err
			}

			var vulnerabilities ent.Vulnerabilities
			for _, results := range report.Results {
				for _, rawVulnerability := range results.Vulnerabilities {
					vulnerabilities = append(vulnerabilities, &ent.Vulnerability{
						VulnerabilityID:         rawVulnerability.VulnerabilityID,
						VulnerabilityURLDetails: rawVulnerability.PrimaryURL,
						InstalledVersion:        rawVulnerability.InstalledVersion,
						FixedVersion:            rawVulnerability.FixedVersion,
						Status:                  vulnerability.Status(rawVulnerability.Status),
						Severity:                vulnerability.Severity(rawVulnerability.Severity),
						Title:                   rawVulnerability.Title,
					})
				}
			}

			if err = w.manifestRepository.CreateVulnerabilitiesInBulkAndMarkAsScanned(grpCtx, vulnerabilities, manifest); err != nil {
				return err
			}

			fmt.Printf("Scan for %s finished\n", opts.ScanOptions.Target)

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func getOpts(baseUrl string, organizationName string, registryName string, imageName string, digest string) (*flag.Options, error) {
	var dbRepositories []name.Reference
	for _, repo := range []string{db.DefaultGCRRepository, db.DefaultGHCRRepository} {
		ref, err := parseRepository(repo, db.SchemaVersion)
		if err != nil {
			return nil, err
		}
		dbRepositories = append(dbRepositories, ref)
	}

	targetName := fmt.Sprintf("%s/%s/%s/%s@%s", baseUrl, organizationName, registryName, imageName, digest)

	return &flag.Options{
		DBOptions: flag.DBOptions{
			DBRepositories: dbRepositories,
		},
		ScanOptions: flag.ScanOptions{
			Target: targetName,
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
