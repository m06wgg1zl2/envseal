# envseal

A tool to encrypt and version-control `.env` files using age encryption with team key management support.

---

## Installation

```bash
go install github.com/yourusername/envseal@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/yourusername/envseal/releases).

---

## Usage

Initialize a new envseal configuration in your project:

```bash
envseal init
```

Encrypt your `.env` file before committing:

```bash
envseal seal .env
```

Decrypt a sealed file locally:

```bash
envseal open .env.sealed
```

Add a teammate's public key to allow them to decrypt:

```bash
envseal keys add age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8z
```

Commit the sealed file safely to version control:

```bash
git add .env.sealed envseal.toml
git commit -m "chore: update sealed env"
```

> **Note:** Never commit your plain `.env` file. Add it to `.gitignore`.

---

## How It Works

`envseal` uses [age](https://github.com/FiloSottile/age) encryption to seal your environment files. Public keys for all authorized team members are stored in `envseal.toml`, which is safe to commit. Each team member uses their own private key to decrypt.

---

## License

MIT © [yourusername](https://github.com/yourusername)