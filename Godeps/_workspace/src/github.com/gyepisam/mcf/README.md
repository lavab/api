#Description

mcf is a Go library for creating, verifying, upgrading and managing a variety of hashed password schemes.

mcf provides a simple API for applications to use a variety of password
hashing schemes, including bcrypt, scrypt, and pbkdf2 as well a management
mechanism to easily and transparently set the default password
scheme, change schemes, or change scheme parameters such as work factors,
salt length, key length without rewriting the application.

In addition to the separation of usage and policy model of the UNIX password
mechanism, mcf borrows the /etc/passwd format and generates passwords
in the [Modular Crypt Format \(MCF\)] (http://pythonhosted.org/passlib/modular_crypt_format.html).
,allows for the simultaneous existence of multiple types and generations of passwords.
As a text format, it provides for easy database storage and subsequent verification.

Any application would benefit from the simplicity, ease and secure
defaults of this package. Applications and web sites that need to support
multiple password hashing mechanisms and/or need to allow for different
upgrade policies would find it especially useful.

#Install

mcf can be installed with the command:

    go get github.com/gyepisam/mcf

#Example

Here is a simple example on how to use mcf to generate pbkdf2 hashed passwords.
Note that the same example would work for any other scheme and multiple schemes
can exist concurrently.

    import "github.com/gyepisam/mcf"
    import "github.com/gyepisam/mcf/pbkdf2"


To create a password, assuming you have the `username` and `plaintext` variables from a user signup session:

    passwd, err := mcf.Create(plaintext)
    // error handling elided
    // Store passwd in database

To verify a password, assuming you have the `username` and `plaintext` variables from a user login session
and have retrieved the stored password from the database into the `passwd` variable:

    isValid, err := mcf.Verify(plaintext, passwd)
    // error handling elided.

The verification can be combined with automated password upgrades with code like this:

    isValid, err := mcf.Verify(plaintext, passwd)
    // error handling elided.
    
    if !isValid {
        // fail -- failed password verification
        return false
    }

    isCurrent, err := mcf.IsCurrent(passwd) 
    // error handling elided

    if !isCurrent {
      newPasswd, err := mcf.Create(plaintext)
      // error handling elided
      // Update passwd in database
    }

    // success!

With this method, user passwords will automatically upgrade to use any new schemes or stronger work factors
set by policy through the management API.

#Author 

mcf is written by Gyepi Sam <self-github@gyepi.com> and is released under a BSD license.

I am welcome any and all feedback on this software.

-Gyepi Sam
