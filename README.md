# MyTime (Go Version)

The motivation behind this project (as with all other versions of MyTime) has been to learn [Go](https://go.dev/). Unlike the [Rust](https://github.com/francescarpi/mytime-gui) version, which is more stable, I do not recommend using this version in a production environment. I am sharing it purely for educational purposes. (Although I must admit that I am using it myself, and I like it more than the Rust/Tauri version).

This version is compatible with the same SQLite database as the Rust version, but it lacks the configuration features.

## Additional Notes

- **Why Go?**: Go was chosen for its simplicity, performance, and strong concurrency model, making it an excellent language for learning and building efficient tools.
- **Future Plans**: While this version is currently experimental, future updates may include configuration support and additional features to bring it closer to production readiness.

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/francescarpi/mytime-go.git
   cd mytime-go
   ```

2. Build the project:
   ```bash
   go build -o bin/mytime ./cmd/ui
   ```

3. Run the application:
   ```bash
   ./bin/mytime
   ```

4. Ensure you have an existing SQLite database compatible with the Rust version.

Let me know if you'd like me to expand on any specific section or add more details!

