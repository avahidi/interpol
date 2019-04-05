

.. image:: logo.png
   :align: center

Interpol
========

**Interpol** is a minimal `string interpolation <https://en.wikipedia.org/wiki/String_interpolation>`_
library written in Go. It can be used to generate a series of strings from a set of rules.
This is useful for example for people doing penetration testing or fuzzing.


**Police** is a command-line interface for Interpol. It is not as powerful as embedding Interpol in your
own application (which allows you to create custom interpolators and modifiers) but still very handy if you are a
CLI type of person.

You can install Police from the `Snap store <https://snapcraft.io/police>`_ ::

    $ sudo snap install police

To build Police from source, install the Go compiler then execute this::

    $ go get -u bitbucket.org/vahidi/interpol/cmd/police/...


Usage example
-------------

Consider the following problem: You have forgotten your password to the company mainframe.
You do however remember that the password had the following format::

    <one of the Friend characters> <a digit> <a currency sign>

Since this is something that can be defined as a bunch of rules, we can use police to generate all possible combinations::

    # 'friends.txt' is a file containing one friends characters per line
    $ police "{{file filename='friends.txt'}}{{counter min=0 max=9}}{{set data='£$¥€'}}"

    Rachel0£
    Monica0£
    Phoebe0£
    . . .
    Joey9€
    Chandler9€
    Gunther9€

You may now use these candidates with a password recovery tool to find your lost password in no time.
There are of course other tools for this particular use case, but I believe few have the flexibility of Interpol/Police.


Interpolators
-------------

An interpolation has the following syntax::

    {{type parameter1=value1 parameter2=value2 ... }}

For example::

    {{counter min=1 max=10 step=3}}

The following interpolators are currently available::

    {{counter [min=0] [max=10] [step=1] [format="%d] }}
    {{random [min=0] [max=100] [count=5] [format="%d"] }}
    {{file filename="somefile" [count=-1] [mode=linear] }}
    {{set data="some input" [sep=""] [count=-1] [mode=linear] }}
    {{copy from="name of another interpolator" }}

Where

- [parameter=value] indicates an optional parameter, value is the default value
- valid values for mode are: linear, random or perm
- format is standard Go fmt.Printf() format string
- copy repeats the value of another interpolator


Copying
~~~~~~~

Interpolators may be given a name. This is needed when using copy::

    {{counter name=mycounter}} {{copy from=mycounter}}

This will yield "0 0", "1 1", and so on.


Modifiers
~~~~~~~~~

Interpolators can also have a *modifier*, which changes their output.
Currently the following modifiers exist:

- *toupper*: make all characters upper case
- *tolower*: make all characters lower case
- *capitalize*: capitalize each word
- *1337*: leet speak modifier (random upper/lower case)

For example, the following will yield "Yes", "No" and "Maybe"::

    {{set data="YES,no,mayBE" sep="," modifier=capitalize}}


Examples
--------

The folder examples/ contains the following samples:

- **hackernews** - download 3 random HN comments from firebase
- **nena** - demonstrates use of copy
- **hodor** - as the name clearly implies this one teaches you to create custom interpolators
- **discordia** - demonstrates use of custom modifiers


License
-------

This library is licensed under the GNU GENERAL PUBLIC LICENSE, version 2 (GPLv2).

See the file LICENSE for more information.

