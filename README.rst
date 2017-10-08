

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

Assume you have been given the task of finding employees who use weak passwords.
You are given a file containing all 100 usernames and another file containing
100 weak passwords.


You can instruct Interpol to use these files as input like so::

    // note: error handling not shown below
    ip := interpol.New()
    user, err := ip.Add("{{file filename=usernames.txt}}")
    password, err := ip.Add("{{file filename=weakpasswords.txt}}")
    for {
        if checkCredentials( user.String(), password.String()) {
            report(user.String() )
        }
        if ! ip.Next() {
            break
        }
    }

Hence interpol will generate all 10.000 valid outputs (username/password pairs)
given the set of rules (in this case the input files). But you could do all
this with a simple for-loop so lets try something more interesting.

Assume you suspect user "joe" is using a password that is a combination of
a weak password plus two additional characters, the first one being a number
and the second a currency sign. You can now narrow your search by doing this::

    // again, error checks omitted
    user, err := ip.Add("joe")
    password, err := ip.Add("{{file filename=weakpasswords.txt}}{{counter min=0 max=9}}{{set data=$£€}}")

The first string is static, the second one however has 3 interpolated elements.
This configuration will generate 1 * 100 * 10 * 3 = 3000 username/password pairs.


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

- *name* is used to name an element (used by copy)
- *modifier* defines an output modifier

Currently the following modifiers exist:

- *toupper*: make all characters upper case
- *tolower*: make all characters lower case
- *1337*: leet speak modifier (random upper/lower case)

Note that you can create your own interpolators and modifiers. 
See the examples for more information.

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


