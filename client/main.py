from diffiehellman.diffiehellman import DiffieHellman

alice = DiffieHellman()
bob = DiffieHellman()

alice.generate_public_key()    # automatically generates private key
bob.generate_public_key()

alice.generate_shared_secret(bob.public_key, echo_return_key=True)
bob.generate_shared_secret(alice.public_key, echo_return_key=True)

a = alice.shared_key
print(a)
b = bob.shared_key
print(b)
