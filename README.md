# 403Unlocker-CLI

403Unlocker-CLI is a versatile command-line tool designed to bypass 403 restrictions effectively. It provides subcommands to handle DNS resolution, DNS server selection, and Docker image proxy discovery.

## Features
- **Check**: Test if a specific URL can be resolved using a custom DNS server.
- **DNS**: Find the most responsive DNS server from a list of custom DNS options.
- **Docker**: Identify the best Docker image proxy to bypass network restrictions.

---

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/403unlocker/403Unlocker-cli.git
   cd 403Unlocker-cli
   ```
2. Install the project:
   ```bash
   make install
   ```

3. Run the project:
   ```bash
   403unlocker --help
   ```
---

## Usage

### General Syntax
```
403unlocker <command> [flags]
```

### Commands

#### 1. Check
Test if a URL can be resolved using a custom DNS server.
```
403unlocker check <URL>
```
Example:
```
403unlocker check "https://pkg.go.dev"
```

#### 2. DNS
Find the fastest DNS sni-proxy among a list of DNS options.
```
403unlocker dns <URL>
```

Example:
```
403unlocker dns "https://packages.gitlab.com/gitlab/gitlab-ce/packages/el/7/gitlab-ce-16.8.0-ce.0.el7.x86_64.rpm/download.rpm"
```

#### 3. Docker
Identify the best Docker image proxy for bypassing network restrictions.
```
403unlocker docker <DOCKER-IMAGE>
```

Example:
```
403unlocker docker "gitlab/gitlab-ce:17.0.0-ce.0"
```


---

## Flags
- `--help`: Display help for any command.

---

## Requirements
- Go 1.18 or higher

---

## Contributing
Contributions are welcome! Feel free to open an issue or submit a pull request.

1. Fork the repository.
2. Create a new branch.
3. Commit your changes.
4. Push the branch and create a PR.

---

## License
This project is licensed under the GPL-3.0 License. See the [LICENSE](https://github.com/403unlocker/403Unlocker-cli/blob/main/LICENSE) file for more information.

---

## Contact
For any questions or feedback, reach out at [borhanisaleh6@gmail.com](mailto:borhanisaleh6@gmail.com).