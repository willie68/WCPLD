# wcupl

Since WinCUPL is no longer running on any of my computers, I wrote a batch, directly starting cupl.exe, with which I can generate the JED file and also run the simulator right away. Personally, though, I like having the sources in one place. I personally find the duplication of the headers of the PLD and SI file and the manual synchronization between the two unpleasant. That's why I came up with the simple WPLD format and wrote the WCUPL tool for it. The WCUPL tool is available in my repo (https://github.com/willie68/WCPLD/releases). This eliminates the tedious header matching and logic and test are in one file. 

The tool itself doesn't do much. Any arguments are passed directly to CUPL. 

- macros with {{}} will be exchange with some other data:
  `{{filename}}`: is the filename without wpld extension
  `{{date}}`: is the actual date in format yyyy-mm-dd (ISO8601)
- it will create a subfolder called wcupl and copy the wpld into it. 
- From this wpld files, `header:` and `pld:` are merged into one #.pld file and `header`: and `simulator:` part into the #.si file. 
- Then cupl is started from the same directory where wcupl is located. Or if present using environment variable `CUPL_HOME` for trying to find cupl.exe. (So sorry, this version is only for windows users)



