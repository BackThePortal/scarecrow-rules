# Scarecrow Bot Rules Manager

Simple Go script that uses the Google Docs API to fetch a document and update the rules of TheTinMen server
automatically.

The program currently supports the following formatting:

- Headings up to level 3 (Discord doesn't support further headings)
- **bold** and _italic_ text
- Bullet lists

In the future, there will be support for:

- Links
- Indentation in lists

Since the program is only meant to output plain text, rich links and mentions (including channel and message links)
aren't possible because the script is completely separated from the Discord API.

## Setup

When you run the program for the first time, you will be asked to give the OAuth2 permission using your Google Account.
Once you complete all your steps, you might be redirected to a page in `localhost` and get a "connection refused" error. This is a bug in Google's API. To continue, you must take the part of the URL after `&code=` until the next `&` and copy it to the terminal.

## Arguments

The syntax is as follows:

```
./main --
```


