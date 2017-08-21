Hugo Redirect
=============

If you have been blogging on wordpress and decided to switch to Hugo but don't really want to move all your wordpress stuffs to Hugo just yet (or never) but looking for a way to have links on your Hugo website to forward to your wordpress blog, this may come in handy.

![gifs.com](https://j.gifs.com/Anw5XO.gif)

# Usage

#### Pre-requisite
Add the following line to the `single.html` file,
```
{{ if .Params.redirectURL }}<meta http-equiv="refresh" content="1; url={{ .Params.redirectURL }}"/>{{ end }}
```

Then run the following to generate redirect pages (this is typically run from your hugo site directory),
```
./hugo-redirect -l content/archive -f wp.xml
```

`-l` Location of the archive

`-f` Wordpress XML file that you get from wordpress blog

Here is how it works:
1. Reads the XML file
2. Creates pages with `redirectURL` and specific `tags` and `categories`

A markdown page created by `hugo-redirect` looks like this,
```
---
title: Monitoring with Nagios
redirectURL: https://mdshaonimran.wordpress.com/2012/11/24/monitoring-with-nagios/
date: 2012-11-23T19:48:14-00:00
tags:
- monitoring
- nagios
categories:
- Monitoring
- Open Source
---
```

## Note
In wordpress Tags and Categories are case-sensetive and sometimes looks like duplicates when moving to other static blogs like Nikola or Hugo. While creating pages for redirect, it doesn't add the duplicates.


3. `redirectURL` parameter creates the following `<meta>` when Hugo generates the HTML file.
4. You need to add the following line to your theme's single.html* file.
```
{{ if .Params.redirectURL }}<meta http-equiv="refresh" content="1; url={{ .Params.redirectURL }}"/>{{ end }}
```
* single.html file name or location may vary based on the theme.

This generates a `<meta>` HTML tag for all the WP blog pages,
```
<meta http-equiv="refresh" content="1; url=https://mdshaonimran.wordpress.com/2012/11/24/monitoring-with-nagios/"/>
```

5. Run the `hugo-redirect` command with location and XML file, it should generate markdown files for all the blogs found in the XML file, except the re-blogs.

Now run `hugo` as you would do for your blog. You should be able to see all the WP blog posts with tags and categories. If clicked on the blogs, you should be redirected to the appropriate WP page.