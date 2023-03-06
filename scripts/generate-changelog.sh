#!/bin/bash

resources=()
data_resources=()
enhancements=()
bugs=()

target_issues=$(sed -e 's/[^0-9]//g' issues_in_release.json )
for issue_number in $target_issues; do
    issue_content=$(gh issue view $issue_number)
    issue_title_content=$(echo "$issue_content" | grep -E "title" | awk -F ": " '{print $NF}')
    issue_title=$(echo $issue_title_content | awk -F ": " '{print $NF}')
    issue_label_content=$(echo "$issue_content" | grep -E "labels")
    issue_label=$(echo $issue_label_content | awk -F " " '{print $NF}')
    
    release_note_content="* \`$issue_title\` (#$issue_number)"
    
    case "$issue_label" in
        "enhancement")
            enhancements+=("$release_note_content")
            ;;
        "resource"|"data resource")
            resources+=("$release_note_content")
            ;;
        "bug")
            bugs+=("$release_note_content")
            ;;
        *)
            echo "$issue_title contains the $issue_label that isn't supported in release notes."
            ;;
    esac
done

printf '%s\n\n' "# <replace with release version> $(date +'%B %d %Y')" >> header.md

if [ ${#enhancements[@]} -gt 0 ]; then
    printf '%s' "### ENHANCEMENTS" >> new_enhancements.md
    printf '\n%s\n' "${enhancements[@]}" >> new_enhancements.md
fi

if [ ${#resources[@]} -gt 0 ]; then
    printf '%s' "### RESOURCES" >> new_resources.md
    printf '\n%s\n' "${resources[@]}" >> new_resources.md
fi

if [ ${#bugs[@]} -gt 0 ]; then
    printf '%s' "### BUG FIXES" >> bug_resolutions.md
    printf '\n%s\n' "${bugs[@]}" >> bug_resolutions.md
fi
    
cat  header.md new_enhancements.md new_resources.md bug_resolutions.md > ../CHANGELOG.md

rm -f header.md new_enhancements.md new_resources.md bug_resolutions.md