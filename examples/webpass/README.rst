"webpass" is a simple web form password brute force utility, where passwords and usernames are given as Interpol expressions. For example::

    webpass --negtext "Invalid" -username "{{file filename=usernames.txt}}" -password "{{file filename=passwords.txt}}" -url "http://localhost/test/login"                                           

It started as an experiment at work, and I felt it would do a nice example on real-world use of Interpol.
