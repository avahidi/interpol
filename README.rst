
Interpol
========

Interpol is an `string interpolation <https://en.wikipedia.org/wiki/String_interpolation>`_
library for people doing security research or penetration testing.

The library is written in the Go programming language and has a very simple yet flexible API.


Background
----------

There is a bit of story behind this library, involving a hacker and a bunch mischievous security researchers.
To make a long story short, someone had set up a phishing page impersonating our company. The login page would post this to the hacker::

    http://company.com.fake.com/login?name=Joe%20Schmuck&email=joes@company.com&password=qwerty


So what if you could build a script that accessed the following URL 1 million times::

    http://company.com.fake.com/login?name=<FirstName>%20<LastName>&email=<UserName>@company.com&password=<Password>

. . . with <FirstName> & <LastName> randomly selected from a list of common Western names and <UserName> & <Password> from lists of common usernames and weak passwords.
And just out of habit, lets pick "User-Agent" from a list of common browser signatures...


And this was when the idea behind this library was born...


Using interpol
--------------

Consider the following example::

    import "bitbucket.org/vahidi/interpol"
    . . .
    ip := interpol.New()
    c, err := ip.Add("n={{counter min=5 max=7}}")
    . . .
    for {
        fmt.Println(c.String() )
        if ! ip.Next() {
            break
        }
    }


Here {{counter min=5 max=7}} is a "counter" interpolator. Running this code will output::

    n=5
    n=6
    n=7

Other common interpolators are "file" and "random". Consider this example::

    // assume domains.txt contains a list of domain names such as microsoft.com and google.com
    c, err := ip.Add("http//{{random min=0 max=999 count=50 format='%03d'}}.{{file filename=domains.txt mode=perm}}")

This could for example generate the following strings::

    http://000.microsoft.com
    http://900.google.com
    http://700.microsoft.com
    ...


License
-------

This library is licensed under the GNU GENERAL PUBLIC LICENSE, version 2 (GPLv2).

See the file COPYING for more information


