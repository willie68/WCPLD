# wcupl

Since WinCUPL is no longer running on any of my computers, I wrote a batch, directly starting cupl.exe, with which I can generate the JED file and also run the simulator right away. Personally, though, I like having the sources in one place. I personally find the duplication of the headers of the PLD and SI file and the manual synchronization between the two unpleasant. That's why I came up with the simple WPLD format and wrote the WCUPL tool for it. The WCUPL tool is available in my repo (https://github.com/willie68/WCPLD/releases). This eliminates the tedious header matching and logic and test are in one file. 

The tool itself doesn't do much. Any arguments are passed directly to CUPL. 

- macros surrounded by `{{}}` will be exchange with some other data:
  `{{filename}}`: is the filename without wpld extension
  `{{date}}`: is the actual date in format yyyy-mm-dd (ISO8601)
- it will create a subfolder called wcupl and copy the wpld into it. 
- From this wpld files, `header:` and `pld:` are merged into one #.pld file and `header`: and `simulator:` part into the #.si file.  The section names `header, pld, simulator` are case insensitive, so `header`, `HEADER` or even `hEaDeR` means the same.
- Then cupl is started from the same directory where wcupl is located. Or, if present, using environment variable `CUPL_HOME` for trying to find cupl.exe. (So sorry, this version is only for windows users)



Example wpld file:

```wpld
HEADER:
Name     {{filename}} ;
PartNo   01 ;
Date     {{date}} ;
Revision 02 ;
Designer wkla ;
Company  nn ;
Assembly None ;
Location  ;
Device   G16V8 ;

PLD:
/* *************** INPUT PINS *********************/
PIN 1   =  A12; 
PIN 2   =  A13;
PIN 3   =  A14;
PIN 4   =  A15;
PIN 9   =  PHI2;

/* *************** OUTPUT PINS *********************/
PIN 12   =  CSRAM;
PIN 13   =  CSHIROM;
PIN 15   =  CSIO;

/* *************** LOGIC *********************/

/* RAM */
CSRAM = A15 # !PHI2;

/* 8kb of ROM */
CSHIROM = !(A15 & A14 & A13);

/* IO */
CSIO= !(A15 & A14 & !A13 & A12);

SIMULATOR:
/* Starting Test description */
ORDER: A15, A14, A13, A12, PHI2, CSRAM, CSHIROM, CSIO; 

VECTORS:
0 X X X 0 H H H 
0 X X X 1 L H H 
1 0 X X X H H H 
1 1 0 0 X H H H 
1 1 0 1 X H H L 
1 1 1 X X H L H 

```

