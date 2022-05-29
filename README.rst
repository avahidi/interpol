

.. image:: logo.png
   :align: center

Interpol
========

**Interpol** is a minimal rule-based `string interpolation <https://en.wikipedia.org/wiki/String_interpolation>`_ library for Go applications.

It can be used in security applications and elsewhere where you want to generate data based on some known format (e.g. passwords that start with a letter and end with a digit). It could for example allow you to create a small corpus instead of doing an exhaustive search.


The library is very simple yet flexible, and even allows you to define your own operators. See the examples/ folder for more information.


Police
------
Police is the Interpol command-line utility. If you don't care about embedding Interpol in your own application and/or adding custom functionality, this is what you should use.



Example
~~~~~~~

Assume you have forgotten your password to the company mainframe.
You do however remember that the password had the following format::


    <one of the Friends characters> <a digit> <a currency sign>


We define these as string interpolations rules:

#. We write a *set* rule for the currency signs: "{{set data='£$¥€'}}"
#. We write a *counter* rule to count from 0 to 9: "{{counter min=0 max=9}}"
#. We create a file named 'friends.txt' with all the characters. Then we create a *file* rule as follows: "{{file filename='friends.txt'}}"

We put all these together as one string and execute it using Police:


    $ police "{{file filename='friends.txt'}}{{counter min=0 max=9}}{{set data='£$¥€'}}"

Which should generate the following output::


    Rachel0£
    Monica0£
    Phoebe0£
    . . .
    Joey9€
    Chandler9€
    Gunther9€

You may now use these candidates with a password recovery tool to find your lost password in no time.


Installing
----------

To install Interpol, first install golang then run this:

    go install github.com/avahidi/interpol@latest

If you are interested in Police, this is what you want instead:

    go install github.com/avahidi/interpol/cmd/police@latest


Documentation
-------------

In our example above the rules were defined as expressions embedded in a string.
Evaluating expressions within strings is often called *string interpolation*,
hence we have chosen to call each rule fragment an "interpolation" and the logic behind it an "interpolator".


An interpolation has the following syntax::

    {{type parameter1=value1 parameter2=value2 ... }}

For example::

    {{counter min=1 max=10 step=3}}


Interpolators
~~~~~~~~~~~~~

The following interpolators are currently available::


    {{counter [min=0] [max=10] [step=1] [format="%d] }}
    {{random [min=0] [max=100] [count=5] [format="%d"] }}
    {{file filename="somefile" [count=-1] [mode=linear] [optional=false] }}
    {{set data="some input" [sep=""] [count=-1] [mode=linear] [optional=false] }}
    {{copy from="others-name" }}

Where

#. [parameter=value] indicates an optional parameter, value is the default value
#. valid values for mode are: linear, random or perm
#. format is standard Go fmt.Printf() format string (which is fairly similar to C format strings)
#. optional means the output is optional (i.e. it may be empty).
#. copy repeats the value of another interpolation (see below)


Copying
~~~~~~~

Interpolators may be given a name. This is needed when using copy:

    "{{counter name=mycounter}} {{copy from=mycounter}}"

This will yield "0 0", "1 1", and so on.


Modifiers
~~~~~~~~~

Interpolators can also have an output *modifier*.
Currently the following modifiers exist:

- *empty*: the empty string "" (ignores input)
- *len*: length of the input (in raw bytes, no fancy UTF-8 support)
- *bitflip*: randomly flip one bit (again, using raw bytes)
- *byteswap*: randomly swap two bytes (raw bytes again)
- *reverse*: reverse (for once, this one supports UTF-8)
- *trim*: trim text (remove space before and after)
- *base64*: base64 encode
- *toupper*: make all characters upper case
- *tolower*: make all characters lower case
- *capitalize*: capitalize each word
- *1337*: leet speak modifier (random upper/lower case)

For example::

    $ police '{{set data="YES,no,mayBE" sep="," modifier=capitalize}}'
    Yes
    No
    Maybe



License
-------

This library is licensed under the GNU GENERAL PUBLIC LICENSE, version 2 (GPLv2).

See the file LICENSE for more information.
