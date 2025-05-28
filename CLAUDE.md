# Agent instructions

Agent guidance for software development.

## Code Standards

### Comments
Add comments that clarify non-obvious code or mark section boundaries. Include doc comments for types/functions when idiomatic.

### Documentation
Propose new docs when important context would be lost. Consider language-specific docs when idiomatic.

### Git Commits
Follow these rules for commit messages:
1. Subject: 50 chars max, capitalized, imperative mood, no period
2. Blank line between subject and body
3. Body: wrap at 72 chars, explain why not how
4. Test: "If applied, this commit will [subject]"
5. Include `Co-Authored-By: Claude <noreply@anthropic.com>` if appropriate
6. No emojis or "Generated with" attribution text

Focus on context and motivation. Make atomic commits for single logical changes.

### White space
ALWAYS avoid trailing white space. Remove it if detected.

### Test first development (TFD)

Use TFD or TDD, unless impractical.  Tests must be falsifiable: we must see
tests fail before we make the code changes for tests to pass.

### Text file line length

Wrap lines for text and Markdown files at 80 characters for readability.

## Agent Tools

If desired CLI tools aren't available, suggest a way to `brew install` them.
