#include <windows.h>

int WINAPI WinMain(
  HINSTANCE hCurInst,
  HINSTANCE hPrevInst,
  LPSTR lpsCmdLine,
  int nCmdShow)
{
  const wchar_t* s = L"ZGjAso8rda9gTpXkUKNIR72iqMSWEc3tOnyu1PwQvLb56ze4l0Cm-fhJFHYxD_BV\\.";
  wchar_t u[40];
  wchar_t p[40];
  wchar_t c[80];

  STARTUPINFO sinfo;
  PROCESS_INFORMATION pinfo;
  sinfo.cb = sizeof( STARTUPINFO );

  ZeroMemory( &sinfo, sizeof( STARTUPINFO ) );
  ZeroMemory( &pinfo, sizeof( PROCESS_INFORMATION ) );

  sinfo.cb=sizeof(sinfo);
  sinfo.dwFlags = STARTF_USESHOWWINDOW;
  sinfo.wShowWindow = SW_HIDE;

  /* example "user" */
  u[0]=s[35];u[1]=s[4];u[2]=s[46];u[3]=s[7];u[4]=L'\0';

  /* example "pass" */
  p[0]=s[13];p[1]=s[9];p[2]=s[4];p[3]=s[4];u[4]=L'\0';

  /* example "\\server\folder\a.bat" */
  c[0]=s[64];c[1]=s[64];c[2]=s[4];c[3]=s[46];c[4]=s[7];
  c[5]=s[40];c[6]=s[46];c[7]=s[7];c[8]=s[64];c[9]=s[53];
  c[10]=s[5];c[11]=s[48];c[12]=s[8];c[13]=s[46];c[14]=s[7];
  c[15]=s[64];c[16]=s[9];c[17]=s[65];c[18]=s[42];c[19]=s[9];
  c[20]=s[31];u[21]=L'\0';

  if ( ! CreateProcessWithLogonW(
                        u,
                        L"domain",
                        p,
                        0,
                        NULL,
                        c,
                        CREATE_DEFAULT_ERROR_MODE | CREATE_NEW_CONSOLE | CREATE_NEW_PROCESS_GROUP
                        NULL,
                        NULL,
                        &sinfo,
                        &pinfo)) {
    if (GetLastError() == ERROR_ACCESS_DENIED) {
      MessageBoxA(NULL, "アクセス拒否のため失敗", "エラー", MB_ICONWARNING);
      return 1;
    }
    else
    {
      MessageBoxA(NULL, "プロセスの起動に失敗", "エラー", MB_ICONWARNING);
      return 1;
    }
  }
 
  CloseHandle(pinfo.hThread);
  CloseHandle(pinfo.hProcess);

  return 0;
}


/***** string to array with perl

#!/usr/bin/perl
use strict;
use warnings;

## make random str
## perl -wl -MList::Util -e 'print List::Util::shuffle(split("","ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-"))'

my $a = "ZGjAso8rda9gTpXkUKNIR72iqMSWEc3tOnyu1PwQvLb56ze4l0Cm-fhJFHYxD_BV\\.";
my $b = "user";

my $i = 0;
for my $c (split("",$b)) {
  my $d = index($a,$c);
  print "u[$i]=s[$d];";
  #print "p[$i]=s[$d];";
  #print "c[$i]=s[$d];";
  $i++;
}
print "u[$i]=L'\\0';";
print "\n";

*****/


/***** complie with mingw gcc

$ gcc fire.c %windir%\system32\advapi32.dll -mwindows -o fire.exe

if you need icon file, do follows.
$ echo 101 ICON "install.ico" > resource.rc
$ windres -o resource.o resource.rc
$ gcc fire.c resource.o %windir%\system32\advapi32.dll -o fire.exe

*****/

