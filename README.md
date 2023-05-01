# Bulk SMS sender

The code is used to send SMS in bulk through MessageBird.

To use it:

- Copy `config.yml.template` to `config.yml`
- Customize the `message` and the `text`. The code now expect a single `%s` in the `message` that will be interpolated with the first column of `names.csv`
- Create a `names.csv` file with 3 columns:
    - The first column should be the first name
    - The second, the last name (not used)
    - The third, the phone number it should be sent to
    - Set your MessageBird API into the `MESSAGEBIRD_API` environment variable and then run with
        - `go run main.go`





