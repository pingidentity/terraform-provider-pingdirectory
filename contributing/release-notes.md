# Generating Release Notes
### *This documented process is to only be used by maintainers of this repository!*
## Prerequisites
- [terraform-provider-pingdirectory](https://github.com/pingidentity/terraform-provider-pingdirectory) repository cloned
- A GitHub Access Token. You can find instructions on how to do this [here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-personal-access-token-classic).
  - **The access token will only need repo, workflow, and read:org permissions**.
- GitHub CLI Tool installed
  - This can be done using `brew` on Mac using the command `brew install gh`. For other OS, use the instructions listed [here](https://cli.github.com/manual/installation).
  
## Important Note
- The Pull Request to be included will need to be related to a specific issue branch (Example: 1-test-issue-to-be-resolved)
  - It is important that the prefixed number ID corresponds to the Issue Number
- Steps to link issue to a branch:
  - Navigate to the GitHub issue on the repository
  - On the right pane under the *Development* section, select *Create a branch*
  - The *generate-changelog.sh* script will look for the prefixed number for getting issue information from the Pull Request provided in the *pull_requests_in_release.json* file

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
- Edit the **pull_requests_in_release.json** file to add GitHub Pull Requests you would like to include in the Changelog for Release Notes (example below).
```
   {
     "pull_requests" : [
       "1",
       "2" 
     ]
   }
```
- Run **generate-changelog.sh**
- Edit the generated CHANGELOG.md:
  - Add release version desired in header
  - Resolve any issue(s) in the *TODO - Multiple issue categories* section
    - Issues are placed in the TODO section due to having multiple issue categories assigned. 
    - It is up to the contributor(s) to determine placement for said issue(s) in CHANGELOG.md before committing and merging
- Push CHANGELOG.md to GitHub and create a PR request to be reviewed

*The generate changelog script can be edited to contain more issue types*