

.. image:: logo.png
   :align: center

Interpol
========

Interpol is a minimal `string interpolation <https://en.wikipedia.org/wiki/String_interpolation>`_
library written in golang.
It can be used to generate a series of strings from a set of rules.
This is useful for example for people doing penetration testing or fuzzing.


Using Interpol
--------------

Assume you have been given the task of finding employees who use a weak password.
You are given a file containing all 100 usernames and another file containing
100 weak passwords. You can instruct Interpol to use these files as input like so::

    import "bitbucket.org/vahidi/interpol"

    // ...

    ip := interpol.New()
    // error checks not shown below.
    user, err := ip.Add("{{file filename=usernames.txt}}")
    password, err := ip.Add("{{file filename=weakpasswords.txt}}")

This creates two objects representing the user name and password.
You can now iterate over all possible values::

    for {
        if checkCredentials( user.String(), password.String()) {
            report(user.String() )
        }
        if ! ip.Next() {
            break
        }
    }

Note that this will result in 100 * 100 = 10.000 username/password pairs.

But you probably don't need a library to do that so lets try something more
interesting.
Assume you suspect user "joe" is using a password that is a combination of
a weak password plus two additional characters, the first one being a number
and the second a currency sign. You can now specify your search by doing this::

    // again, error checks omitted
    user, err := ip.Add("joe")
    password, err := ip.Add("{{file filename=weakpasswords.txt}}{{counter min=0 max=9}}{{set data=$£€}}")

The first string is static, the second one however has 3 interpolated elements.
This configuration will generate 1 * 100 * 10 * 3 = 3000 pairs.


Interpolators
-------------

An interpolation has the following format::

    {{type parameter1=value1 parameter2=value2 ... }}

With the following types and parameters currently implemented:

- **counter**: min, max, step, format
- **random**: min, max, count, format
- **file**: filename, count, mode
- **set**: data, sep, count, mode
- **copy**: from

Where 

- *mode* is any of linear, random or perm
- *format* is standard printf format string (e.g. "0x%08X")
- *copy* repeats the value of another interpolator. target must have a name

Furthermore, all interpolators can include the following optional parameters:

- *name* is used to name an element. This is needed if you want to copy it later
- *modifier* defined an output modifier, such as 'tolower'


More examples
-------------

The folder examples/ contains the following samples:

- **rng** - generate pseudo-random between 0000 and 9999
- **hackernews** - download 3 random HN comments from firebase
- **password** - variation of the example shown above
- **nena** - demonstrates use of copy
- **hodor** - as the name clearly implies this one teaches you to create custom interpolators
- **discordia** - demonstrates use of custom modifiers
- **pocli** - interpol command line tool

License
-------

This library is licensed under the GNU GENERAL PUBLIC LICENSE, version 2 (GPLv2).

See the file LICENSE for more information


