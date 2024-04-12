# Agar.io Clone

Welcome to the GitHub repository for my personal project, a clone of the popular game [Agar.io](https://agar.io/). This project aims to replicate the addictive gameplay of Agar.io while exploring and implementing a custom binary protocol for efficient data serialization and deserialization between JavaScript and Golang.

## Features

- **Custom Binary Protocol:** Utilizes a custom-built binary protocol for communication between the client and server, ensuring fast and efficient data handling.
- **Cross-Language Serialization:** Features custom serialization and deserialization techniques that work seamlessly between JavaScript (client-side) and Golang (server-side).
- **Advanced Game Logic:** A significant portion of the game logic has been implemented, closely mimicking the original Agar.io gameplay mechanics.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:
- [Node.js](https://nodejs.org/)
- [Go](https://golang.org/)

### Installation

1. **Clone the repository:**

```bash
git clone https://github.com/yourusername/agar-clone.git
cd agar-clone
```

2. **Install JavaScript dependencies:**

```bash
npm install
```

3. **Run the server:**

```bash
go run server.go
```

4. **Start the client:**

``bash
npm start
``

This will launch the game client in your default web browser.

## Architecture

This project is divided into two main components:

- **Client:** Written in JavaScript, the client handles user interactions, rendering, and communication with the server.
- **Server:** The server, written in Golang, manages game state, player interactions, and executes game logic.

## Contributing

Contributions are welcome! If you have suggestions for improvements or bug fixes, please open an issue or submit a pull request.

## License

This project is open source and available under the [MIT License](LICENSE.md).

## Contact

If you have any questions or want to get in touch, please send an email to diyor_khaydarov@iota.uz.
