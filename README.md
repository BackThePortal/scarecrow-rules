# Scarecrow Bot Rules Manager

Simple Go script that uses the Google Docs API to fetch a document and update the rules of TheTinMen server automatically.

The program currently supports the following formatting:

- Headings up to level 3 (Discord doesn't support further headings)
- **bold** and _italic_ text
- Bullet lists

In the future, there will be support for:

- Links
- Indentation in lists

Since the program is only meant to output plain text, rich links and mentions (including channel and message links) aren't possible because the script is completely separated from the Discord API.

## Arguments

The syntax is as follows:

```
./main --
```

- 