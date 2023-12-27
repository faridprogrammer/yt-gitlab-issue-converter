## Overview

This is a straightforward tool designed to facilitate the conversion of YouTrack issue export files into GitLab-compatible issue format.

## How to Use

1. **Export YouTrack Issues:**
   Obtain an export file from YouTrack containing your issues.

2. **Prepare Export File:**
   Remove unnecessary columns from the export file. Ensure that the columns are in the following order:
   `id, project, summary, reporter, type, state, assignee, description`

3. **Create Configuration File:**
   Generate a `config.json` file at your preferred location. Specify the file paths and mappings for users, issue types, and states.

   Example `config.json`:
   ```json
   {
       "SrcPath": "",
       "DestPath": "",
       "YoutrackToGitlabUser": {
           "YoutrackUser": "GitlabUser"
       },
       "YoutrackToGitlabType": {
           "YoutrackType": "GitlabType"
       },
       "YoutrackToGitlabState": {
           "YoutrackState": "GitlabState"
       }
   }
   ```
4. **Run the tool:**
Execute the tool by providing the path to the configuration file.

    `.\main.exe -c path/to/config.json`

    Adjust the paths and mappings according to your usage specifications.
