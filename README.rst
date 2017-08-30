

.. image:: logo.svg
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
You are given a file containing all 1000 usernames and another file containing
1000 weak password. You can instruct Interpol to use these files like so::

    import "bitbucket.org/vahidi/interpol"
    
    // ...
    
    ip := interpol.New()
    // error checks not shown below.
    user, err := ip.Add("{{file filename=usernames.txt}}")
    password, err := ip.Add("{{file filename=weakpasswords.txt}}")

This creates two objects representing the user name and password.
You can now iterate over all possible values::

    for {
        if checkCredential( user.String(), password.String()) {
            report(user.String() )
        }
        if ! ip.Next() {
            break
        }
    }

Note that this will result in 100 * 100 = 10.000 username/password pairs.
But you probably don't need a library to do that so lets try something more 
interesting...

Assume you suspect user "joe" is using a password that is a combination of 
a weak password plus two additional characters, the first one being a number
and the second one '$'. You can now narrow down your search by doing this::

    // again, error checks omitted
    user, err := ip.Add("joe")
    password, err := ip.Add("{{file filename=weakpasswords.txt}}{{counter min=0 max=9}}$")

The first string is static, the second one however has 1 static and 2 interpolated elements.
This configuration will generate only 100 *10 = 1000 pairs.


Interpolators
-------------

Currently the following "interpolators" are supported:

 - static text (no interpolation)
 - counter
 - random
 - file

Each support a different number of parameters. 
See the examples for more information.


More examples
-------------

The examples/ folder contains the following samples:

 - rng - generate pseudorandom between 0000 and 9999
 - hackernews - download 3 random HN comments from firebase
 - password - the example shown above


License
-------

This library is licensed under the GNU GENERAL PUBLIC LICENSE, version 2 (GPLv2).

See the file LICENSE for more information


