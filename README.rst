.. image:: logo.png
   :align: center

Interpol
========

**Interpol** is a minimal, rule-based `string interpolation <https://en.wikipedia.org/wiki/String_interpolation>`_ library for Go. It is designed for generating data that follows a specific format, such as test data, password corpuses, or any structured text.

A command-line utility, `police`, is also provided for standalone use.

Example
-------

Assume you have forgotten your password for the company mainframe. You do however remember that the password had the following structure: `<a character from the show Friends><a digit><a currency symbol>`.

Using Interpol, we can generate a list of all such combinations:

1.  Create a file named `friends.txt` containing the characters' names.
2.  Run `police` with the following rules:

.. code-block:: bash

    $ police "{{file filename='friends.txt'}}{{counter min=0 max=9}}{{set data='£$¥€'}}"

The generated output can be used with a password recovery tool such as `john` to quickly find your lost password.

.. code-block::
    
    Rachel0£
    Monica0£
    Phoebe0£
    ...
    Joey9€
    Chandler9€
    Gunther9€


Example (library)
-----------------

Interpol can be integrated into Go applications:

.. code-block:: go

    package main

    import (
        "fmt"
        "log"
        "github.com/avahidi/interpol"
    )

    func main() {
        rules := "{{set data='Hello,Goodbye' sep=','}}, {{set data='World,friends' sep=','}}!"

        ip := interpol.New()
        output, err := ip.Add(rules)
        if err != nil {
            log.Fatalf("Failed to create interpolator: %v", err)
        }

        for ip.Next() {
            fmt.Println(output.String())
        }
    }



This will produce:

.. code-block::

    Hello, World!
    Goodbye, World!
    Hello, friends!
    Goodbye, friends!


The library allows you to define you own operators. See the examples/ folder for more information.



Installation
------------

This will install the `police` command-line utility to ~/go/bin :

.. code-block:: bash

    go install github.com/avahidi/interpol/cmd/police@latest


Documentation
-------------

Interpol follows rules defined as expressions embedded in a string. Evaluating expressions within strings is often called string interpolation, hence we have chosen to call each rule fragment an "interpolation" and the logic behind it an "interpolator".

Syntax
~~~~~~

An interpolation has the following syntax:

.. code-block::

    {{type parameter1=value1 parameter2=value2 ... }}

For example:

.. code-block::

    {{counter min=1 max=10 step=3}}


Interpolators
~~~~~~~~~~~~~

The following interpolators are currently available:

.. code-block::

    ┌──────────────┬──────────────────────────────────────────────────────┐                  
    │ Interpolator │ Description                                          │                  
    ├──────────────┼──────────────────────────────────────────────────────┤                  
    │ counter      │ A sequence of numbers                                │
    │ random       │ A set of random numbers within a given range.        │
    │ file         │ Lines from a file.                                   │
    │ set          │ A set of values in a set                             │
    │ copy         │ Output of another interpolator                       │
    └──────────────┴──────────────────────────────────────────────────────┘ 
  
Each interpolator has a number of parameters, some of which are optional and some have default values:
  
  

.. code-block::

    ┌──────────────┬───────────────┬─────────────────────────────────────────────────┐                  
    │ Interpolator │ Mandatory     │ Optional                                        │ 
    ├──────────────┼───────────────┼─────────────────────────────────────────────────┤                  
    │ counter      │               │ min=0, max=10, step=1, format="%d               │
    │ random       │               │ min=0, max=10, count=5, format="%d              │
    │ file         │ filename      │ count=-1, mode=linear, optional=false           │
    │ set          │ data          │ sep="", count=-1, mode=linear, optional=false   │
    │ copy         │ from          │                                                 │
    └──────────────┴───────────────┴─────────────────────────────────────────────────┘ 
                                                                                                           
                                                                                    

Notes:

- `format` uses the standard Go `fmt.Printf()` format string.
- `optional=true` allows the interpolation to produce an empty output.

Copying                                                                                                 
~~~~~~~                                                                                                 
Interpolators can be given a `name` attribute. This is required when using the `copy` interpolator to repeat the value of another interpolation.                                              
                                                                                                        
.. code-block::                                                                                         

    "{{counter name=mycounter}} {{copy from=mycounter}}"

This will yield "0 0", "1 1", "2 2", and so on.                                                         

Modifiers
~~~~~~~~~

Interpolators can have an output `modifier` to transform the generated value.

For example:

.. code-block:: bash

    $ police '{{set data="YES,no,mayBE" sep="," modifier=capitalize}}'
    Yes
    No
    Maybe

The following modifiers are available:

.. code-block:: bash

    ┌────────────┬──────────────────────────────────────────────────────────────────────────────┐                             
    │ Modifier   │ Description                                                                  │                             
    ├────────────┼──────────────────────────────────────────────────────────────────────────────┤                         
    │ empty      │ Returns an empty string, ignoring the input.                                 │                             
    │ len        │ Returns the length of the input in bytes.                                    │                             
    │ bitflip    │ Randomly flips one bit in the byte representation of the input.              │                             
    │ byteswap   │ Randomly swaps two bytes in the input.                                       │                             
    │ reverse    │ Reverses the input string (UTF-8 aware).                                     │                             
    │ trim       │ Removes leading and trailing whitespace.                                     │                             
    │ base64     │ Base64-encodes the input.                                                    │                             
    │ toupper    │ Converts the input to upper case.                                            │                             
    │ tolower    │ Converts the input to lower case.                                            │                             
    │ capitalize │ Capitalizes each word in the input.                                          │                             
    │ 1337       │ Applies "l33t speak" character substitutions (e.g., `e` -> `3`, `o` -> `0`). │                             
    └────────────┴──────────────────────────────────────────────────────────────────────────────┘      

License
-------

This library is licensed under the GNU GENERAL PUBLIC LICENSE, version 2 (GPLv2).

See the file LICENSE for more information.
