import os

def create_large_file(file_path, size_gb):
    # Define the size of each chunk in bytes
    chunk_size = 1024 * 1024  # 1 MB

    # Calculate the total size in bytes
    total_size = size_gb * 1024 * 1024 * 1024

    with open(file_path, 'wb') as f:
        # Write chunks until the total size is reached
        for _ in range(total_size // chunk_size):
            f.write(b'\0' * chunk_size)

        # Write any remaining bytes
        remaining_bytes = total_size % chunk_size
        if remaining_bytes:
            f.write(b'\0' * remaining_bytes)

if __name__ == "__main__":
    file_path = "large_file.bin"
    size_gb = 4

    create_large_file(file_path, size_gb)
    print(f"Successfully created a {size_gb}GB file at {file_path}")
