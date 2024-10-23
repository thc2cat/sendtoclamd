# sendtoclamd : cli for sending file to a remote clamav clamd for analysis

## Use case

Produce MD5 hash of dir, compare to an old reference, get new file list, remove duplicates, send the file list to a remote clamav.

```shell

Generate MD5 hash of dir
# /local/bin/gem-md5 -d /var/www/html > var_www_html.0
# sort var_www_html.0 > new
# sort previousMD5 > old

Diff from previous MD5 hash
# diff --changed-group-format='%>' --unchanged-group-format=''    new old > diffs

Select only new files 
# cat diffs  | /local/bin/uniqfields | /local/bin/sendtoclamd

```
