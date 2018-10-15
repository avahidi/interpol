

.. image:: logo.png
   :align: center

Interpol
========

**Interpol** is a minimal `string interpolation <https://en.wikipedia.org/wiki/String_interpolation>`_
library written in golang. It can be used to generate a series of strings from a set of rules.
This is useful for example for people doing penetration testing or fuzzing.


**Police** is the Interpol command line interface. It is not as powerful as embedding Interpol in your
own application (which gives you custom interpolators and modifiers) but still very handy if you are a
CLI type of person.

Usage
-----

Consider the following example: You have forgotten your password for the company mainframe.
You do however remember that the password had the following format::

    <one of the Friend characters> <a digit> <a currency sign>

Assuming the file 'friends.txt' contains name of all friends characters, we can generate all possible combination using three interpolators::

    $ police "{{file filename='friends.txt'}}{{counter min=0 max=9}}{{set data='£$¥€'}}"

    Rachel0£
    Monica0£
    Phoebe0£
    . . .
    Joey9€
    Chandler9€
    Gunther9€

Use these candidates with a password recovery tool to find your lost password in no time.
There are of course other tools for this particular usecase, but I believe few have the flexibility of Interpol/Police.

See examples/hackernews for a similar example with some networking.


Interpolators
=============

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
- *format* is standard Go Printf format string (e.g. "0x%08X")
- *copy* repeats the value of another interpolator. target must have a name

Furthermore, all interpolators can include the following optional parameters:

- *name* is used to name an element (used by copy)
- *modifier* defines an output modifier

Currently the following modifiers exist:

- *toupper*: make all characters upper case
- *tolower*: make all characters lower case
- *capitalize*: capitalize each word
- *1337*: leet speak modifier (random upper/lower case)


Examples
========

The folder examples/ contains the following samples:

- **hackernews** - download 3 random HN comments from firebase
- **nena** - demonstrates use of copy
- **hodor** - as the name clearly implies this one teaches you to create custom interpolators
- **discordia** - demonstrates use of custom modifiers


Building
--------

To build Police from source, install the Go compiler and execute this::

    $ go get -u bitbucket.org/vahidi/interpol/cmd/police/...

License
-------

This library is licensed under the GNU GENERAL PUBLIC LICENSE, version 2 (GPLv2).

See the file LICENSE for more information


