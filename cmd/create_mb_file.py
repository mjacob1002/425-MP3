import random
import os
import string

def generate_random_string(size):
    return ''.join(random.choice(string.ascii_letters) for _ in range(size))

def generate_file(file_path, size_mb):
    # Calculate size in bytes
    size_bytes = size_mb * 1024 * 1024

    # Generate random string
    random_data = generate_random_string(size_bytes)

    # Write to file
    with open(file_path, 'w') as file:
        file.write(random_data)

def gen_rando_file(fpath: str, size: int):

    cmd = f"dd if=/dev/urandom of={fpath} bs=1024*1024 count={size}"
    os.system(cmd)

if __name__ == "__main__":
    # Specify file path and size in megabytes
    num_files = 10
    file_path = "random_file"
    size_mb = 3  # Change this to the desired size

    # Generate the file
    for i in range(num_files):
            creation_path = file_path + str(i) + ".txt"
            gen_rando_file(file_path, size_mb)

