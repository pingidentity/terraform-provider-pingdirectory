#!/bin/bash

bugs=()
conflicts=()
data_resources=()
dependencies=()
documentations=()
enhancements=()
resources=()

target_prs=$(sed -e 's/[^0-9]//g' pull_requests_in_release.json )
pr_content=$(gh pr list)

for target_pr in $target_prs; do
    issue_number=$(echo "$pr_content" | grep $target_pr | awk -F "\t" '{print $3}' | awk -F "-" '{print $1}')
    issue_content=$(gh issue view $issue_number)
    issue_title_content=$(echo "$issue_content" | grep -E "title" | awk -F ": " '{print $NF}')
    issue_title=$(echo $issue_title_content | awk -F ": " '{print $NF}')
    issue_label_content=$(echo "$issue_content" | grep -E "labels")
    issue_label=$(echo $issue_label_content | awk -F ": " '{print $NF}')
    issue_label_count=$(echo "$issue_label" | sed -e 's/[,]//g' | wc -w | sed -e "s/ //g")
    
    if [ $issue_label_count -gt 1 ] ; then
        issue_label="conflict"
    fi

    release_note_content="* \`$issue_title\` (#$issue_number)"

    case "$issue_label" in
        "conflict")
            conflicts+=("$release_note_content")
            ;;
        "dependencies")
            dependencies+=("$release_note_content")
            ;;
        "documentation")
            documentations+=("$release_note_content")
            ;;
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

if [ ${#conflicts[@]} -gt 0 ]; then
    printf '%s' "### CONFLICTS" >> conflicts.md
    printf '\n%s\n' "${conflicts[@]}" >> conflicts.md
fi

if [ ${#documentations[@]} -gt 0 ]; then
    printf '%s' "### DOCUMENTATION UPDATES" >> documentations.md
    printf '\n%s\n' "${documentations[@]}" >> documentations.md
fi

if [ ${#dependencies[@]} -gt 0 ]; then
    printf '%s' "### DEPENDENCIES" >> dependencies.md
    printf '\n%s\n' "${dependencies[@]}" >> dependencies.md
fi

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
    
cat  header.md conflicts.md documentations.md dependencies.md new_enhancements.md new_resources.md bug_resolutions.md > ../CHANGELOG.md

rm -f header.md conflicts.md documentations.md dependencies.md new_enhancements.md new_resources.md bug_resolutions.md