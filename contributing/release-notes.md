# Generating Release Notes
### *This documented process is to only be used by maintainers of this repository!*
## Prerequisites
- [terraform-provider-pingdirectory](https://github.com/pingidentity/terraform-provider-pingdirectory) repository cloned
- A GitHub Access Token. You can find instructions on how to do this [here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-personal-access-token-classic).
  - **The access token will only need repo, workflow, and read:org permissions**.
- GitHub CLI Tool installed
  - This can be done using `brew` on Mac using the command `brew install gh`. For other OS, use the instructions listed [here](https://cli.github.com/manual/installation).
  
## Important Note
Issues will need to be assigned to the desired milestone (relating to a release) for this script to work properly.
  
## Development Steps
- Login to the GitHub via `gh` CLI tool:
  - `gh auth login`
  - Select **Github.com**
  - Select **HTTPS**
  - Type `y` to authenticate with your GitHub credentials
  - Select **Paste Your Access Token**, then Paste your access token from the [Prerequisites section](#prerequisites)
- `git pull origin main`
- `git checkout -b release-notes-<version number>`
- `cd` to scripts directory
- Edit **milestone** variable to match desired GitHub Milestone ID in the **generate-changelog.sh** file 
- Run **generate-changelog.sh**
- Edit the generated CHANGELOG.md:
  - Add release version desired in header
  - Resolve any issue(s) in the *TODO - Multiple issue categories* section
    - Issues are placed in the TODO section due to having multiple issue categories assigned. 
    - It is up to the contributor(s) to determine placement for said issue(s) in CHANGELOG.md before committing and merging
- Push CHANGELOG.md to GitHub and create a PR request to be reviewed

*The generate changelog script can be edited to contain more issue types*