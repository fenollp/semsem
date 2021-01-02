package semsem

import (
	// "bytes"
	// "encoding/gob"
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/exp/apidiff"
	"golang.org/x/tools/go/packages"
	// "gopkg.in/src-d/go-git.v4"
	// "gopkg.in/src-d/go-git.v4/plumbing"
)

// T TODO
type T struct {
	Pkgs    []*packages.Package
	Self    map[string]*packages.Package
	Imports map[string]*packages.Package
	OldObjs map[string]types.Object
	Report  apidiff.Report
}

// X TODO
func X(pkgPath string) (x *T, err error) {
	cfg := packages.Config{
		Mode:  packages.LoadTypes,
		Tests: false,
	}
	var pkgs []*packages.Package
	if pkgs, err = packages.Load(&cfg, pkgPath); err != nil {
		return
	}

	selfPkgs := make(map[string]*packages.Package)
	importPkgs := make(map[string]*packages.Package)
	for _, pkg := range pkgs {
		// skip internal packages since they do not contain public APIs
		if strings.HasSuffix(pkg.PkgPath, "/internal") || strings.Contains(pkg.PkgPath, "/internal/") {
			continue
		}
		selfPkgs[pkg.PkgPath] = pkg
	}
	for _, pkg := range pkgs {
		for _, ipkg := range pkg.Imports {
			if _, ok := selfPkgs[ipkg.PkgPath]; !ok {
				importPkgs[ipkg.PkgPath] = ipkg
			}
		}
	}

	if len(pkgs) != 1 {
		panic(fmt.Sprintf("FIXME: %+v", pkgs))
	}
	// pkg:=types.NewPackage(pkgs[0].PkgPath, pkgs[0].Name)
	pkg := pkgs[0].Types

	var gkp *types.Package
	{
		gkp = pkg
		// var network bytes.Buffer
		// enc := gob.NewEncoder(&network)
		// if err = enc.Encode(pkg); err != nil {
		// 	return
		// }
		// dec := gob.NewDecoder(&network)
		// if err = dec.Decode(&gkp); err != nil {
		// 	return
		// }
	}

	oldobjs := make(map[string]types.Object)
	for _, name := range pkg.Scope().Names() {
		oldobj := pkg.Scope().Lookup(name)
		if !oldobj.Exported() {
			continue
		}
		oldobjs[name] = oldobj
	}

	report := apidiff.Changes(pkg, gkp)
	// https://github.com/golang/exp/blob/eab1b5eb1a030d481efa6d2d1362ece0ebdba8e9/apidiff/apidiff.go#L113

	x = &T{
		Pkgs:    pkgs,
		Self:    selfPkgs,
		Imports: importPkgs,
		OldObjs: oldobjs,
		Report:  report,
	}
	return
}

// type Options struct {
// 	RepoPath       string
// 	OldCommit      string
// 	NewCommit      string
// 	CompareImports bool
// }

// func Run(opts Options) (*Diff, error) {
// 	repo, err := git.PlainOpen(opts.RepoPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open git repo: %w", err)
// 	}

// 	wt, err := repo.Worktree()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get git worktree: %w", err)
// 	}
// 	if stat, err := wt.Status(); err != nil {
// 		return nil, fmt.Errorf("failed to get git status: %w", err)
// 	} else if !stat.IsClean() {
// 		return nil, fmt.Errorf("git tree is dirty")
// 	}

// 	origRef, err := repo.Head()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get current HEAD reference: %w", err)
// 	}

// 	oldHash, newHash, err := getHashes(repo, plumbing.Revision(opts.OldCommit), plumbing.Revision(opts.NewCommit))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to lookup git commit hashes: %w", err)
// 	}

// 	defer func() {
// 		if err := checkoutRef(*wt, *origRef); err != nil {
// 			fmt.Printf("WARNING: failed to checkout your original working commit after diff: %v\n", err)
// 		}
// 	}()

// 	selfOld, importsOld, err := getPackages(*wt, *oldHash)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get packages from old commit %q (%s): %w", opts.OldCommit, oldHash, err)
// 	}

// 	selfNew, importsNew, err := getPackages(*wt, *newHash)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get packages from new commit %q (%s): %w", opts.NewCommit, newHash, err)
// 	}

// 	diff := &Diff{}
// 	diff.selfReports, diff.selfIncompatible = compareChangesAdditionsAndRemovals(selfOld, selfNew)

// 	if opts.CompareImports {
// 		// When comparing imports, we only compare the changes and additions
// 		// between oldPkgs and newPkgs. We ignore removals in newPkgs because
// 		// the removed packages are no longer dependencies and therefore have
// 		// no impact on compatibility of imports.
// 		diff.importsReports, diff.importsIncompatible = compareChangesAndAdditions(importsOld, importsNew)
// 	}

// 	return diff, nil
// }

// func compareChangesAdditionsAndRemovals(oldPkgs, newPkgs map[string]*packages.Package) (map[string]apidiff.Report, bool) {
// 	reports, incompatible := compareChangesAndAdditions(oldPkgs, newPkgs)

// 	// remove packages from oldPkgs that are present in newPkgs. When this loop
// 	// completes, the packages left in oldPkgs are the ones that were removed
// 	// and no longer used in the new commit of this repo.
// 	//
// 	// This is required for the next loop to be able to report correctly on
// 	// removes between the old commit and new commit.
// 	for k := range newPkgs {
// 		delete(oldPkgs, k)
// 	}

// 	for k, oldPackage := range oldPkgs {
// 		report := apidiff.Changes(oldPackage.Types, types.NewPackage(k, oldPackage.Name))
// 		for _, c := range report.Changes {
// 			if !c.Compatible {
// 				incompatible = true
// 			}
// 		}
// 		reports[k] = report
// 	}
// 	return reports, incompatible
// }

// func compareChangesAndAdditions(oldPkgs, newPkgs map[string]*packages.Package) (map[string]apidiff.Report, bool) {
// 	reports := map[string]apidiff.Report{}
// 	incompatible := false
// 	for k, newPackage := range newPkgs {
// 		// if this is a brand new package, use a dummy empty package for
// 		// oldPackage, so that everything in newPackage is reported as new.
// 		oldPackage, ok := oldPkgs[k]
// 		if !ok {
// 			oldPackage = &packages.Package{Types: types.NewPackage(newPackage.PkgPath, newPackage.Name)}
// 		}

// 		report := apidiff.Changes(oldPackage.Types, newPackage.Types)
// 		for _, c := range report.Changes {
// 			if !c.Compatible {
// 				incompatible = true
// 			}
// 		}
// 		reports[k] = report
// 	}
// 	return reports, incompatible
// }

// func getHashes(repo *git.Repository, oldRev, newRev plumbing.Revision) (*plumbing.Hash, *plumbing.Hash, error) {
// 	oldCommitHash, err := repo.ResolveRevision(oldRev)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("could not get hash for %q: %v", oldRev, err)
// 	}

// 	newCommitHash, err := repo.ResolveRevision(newRev)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("could not get hash for %q: %v", newRev, err)
// 	}

// 	return oldCommitHash, newCommitHash, nil
// }

// func getPackages(wt git.Worktree, hash plumbing.Hash) (map[string]*packages.Package, map[string]*packages.Package, error) {
// 	if err := wt.Checkout(&git.CheckoutOptions{
// 		Hash: hash,
// 	}); err != nil {
// 		return nil, nil, err
// 	}

// 	cfg := packages.Config{
// 		Mode:  packages.LoadTypes,
// 		Tests: false,
// 	}
// 	pkgs, err := packages.Load(&cfg, "./...")
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	selfPkgs := make(map[string]*packages.Package)
// 	importPkgs := make(map[string]*packages.Package)
// 	for _, pkg := range pkgs {
// 		// skip internal packages since they do not contain public APIs
// 		if strings.HasSuffix(pkg.PkgPath, "/internal") || strings.Contains(pkg.PkgPath, "/internal/") {
// 			continue
// 		}
// 		selfPkgs[pkg.PkgPath] = pkg
// 	}
// 	for _, pkg := range pkgs {
// 		for _, ipkg := range pkg.Imports {
// 			if _, ok := selfPkgs[ipkg.PkgPath]; !ok {
// 				importPkgs[ipkg.PkgPath] = ipkg
// 			}
// 		}
// 	}

// 	// Reset the worktree. Sometimes loading the packages can cause the
// 	// worktree to become dirty. It seems like this occurs because package
// 	// loading can change go.mod and go.sum.
// 	//
// 	// TODO(joelanford): If go-git starts to support checking out of specific
// 	//   files we can update this to be less aggressive and only checkout
// 	//   go.mod and go.sum instead of resetting the entire tree.
// 	if err := wt.Reset(&git.ResetOptions{
// 		Mode:   git.HardReset,
// 		Commit: hash,
// 	}); err != nil {
// 		return nil, nil, fmt.Errorf("failed to hard reset to %v: %w", hash, err)
// 	}

// 	return selfPkgs, importPkgs, nil
// }

// func checkoutRef(wt git.Worktree, ref plumbing.Reference) (err error) {
// 	if ref.Name() == "HEAD" {
// 		return wt.Checkout(&git.CheckoutOptions{Hash: ref.Hash()})
// 	}
// 	return wt.Checkout(&git.CheckoutOptions{Branch: ref.Name()})
// }
