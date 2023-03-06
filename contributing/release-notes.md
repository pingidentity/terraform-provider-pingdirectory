# Generating Release Notes
## Prerequisites
- [terraform-provider-pingdirectory](https://github.com/pingidentity/terraform-provider-pingdirectory) repository cloned
- A GitHub Access Token. You can find instructions on how to do this [here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-personal-access-token-classic).
  - **The access token will only need repo, workflow, and read:org permissions**.
- GitHub CLI Tool installed
  - This can be done using `brew` on Mac using the command `brew install gh`. For other OS, use the instructions listed [here](https://cli.github.com/manual/installation).
  
## Development Steps
- Login to the GitHub via `gh` CLI tool:
  - `gh auth login`
  - Select **Github.com**
  - Select **HTTPS**
  - Type `y` to authenticate with your GitHub credentials
  - Select **Paste Your Access Token**, then Paste your access token from the [Prerequisites section](#prerequisites)
- `git pull origin main`
- `git checkout -b <your branch name>`
- `cd` to scripts directory
- Edit the **issues_in_release.json** file to add GitHub issues you would like to include in the Changelog for Release Notes (example below).
```
   {
     "issues" : [
       "76",
       "77"   
     ]
   }
```
- Run **generate-changelog.sh**
- Edit CHANGELOG.md to contain release version desired
- Push CHANGELOG.md to GitHub and create a PR request to be reviewed

*The generate changelog script can be edited to contain more issue types*