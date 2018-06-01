package main

import (
  "fmt"

  "github.com/jroimartin/gocui"
)

func handleBranchPress(g *gocui.Gui, v *gocui.View) error {
  branch := getSelectedBranch(v)
  if output, err := gitCheckout(branch.Name, false); err != nil {
    createSimpleConfirmationPanel(g, v, "Error", output)
  }
  return refreshSidePanels(g, v)
}

func handleForceCheckout(g *gocui.Gui, v *gocui.View) error {
  branch := getSelectedBranch(v)
  return createConfirmationPanel(g, v, "Force Checkout Branch", "Are you sure you want force checkout? You will lose all local changes (y/n)", func(g *gocui.Gui, v *gocui.View) error {
    if output, err := gitCheckout(branch.Name, true); err != nil {
      createSimpleConfirmationPanel(g, v, "Error", output)
    }
    return refreshSidePanels(g, v)
  }, nil)
}

func getSelectedBranch(v *gocui.View) Branch {
  lineNumber := getItemPosition(v)
  return state.Branches[lineNumber]
}

// may want to standardise how these select methods work
func handleBranchSelect(g *gocui.Gui, v *gocui.View) error {
  renderString(g, "options", "space: checkout, f: force checkout")
  if len(state.Branches) == 0 {
    return renderString(g, "main", "No branches for this repo")
  }
  go func() {
    lineNumber := getItemPosition(v)
    branch := state.Branches[lineNumber]
    diff, _ := getBranchDiff(branch.Name, branch.BaseBranch)
    renderString(g, "main", diff)
  }()
  return nil
}

// refreshStatus is called at the end of this because that's when we can
// be sure there is a state.Branches array to pick the current branch from
func refreshBranches(g *gocui.Gui) error {
  g.Update(func(g *gocui.Gui) error {
    v, err := g.View("branches")
    if err != nil {
      panic(err)
    }
    state.Branches = getGitBranches()
    v.Clear()
    for _, branch := range state.Branches {
      fmt.Fprintln(v, branch.DisplayString)
    }
    resetOrigin(v)
    return refreshStatus(g)
  })
  return nil
}