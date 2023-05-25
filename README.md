# askimg 🖼️❓

**askimg** answers questions about images using AI

This is a golang wrapper for [andreasjansson/blip-2](https://replicate.com/andreasjansson/blip-2)

## 📦 Installation

You can use the Golang binary to install **askimg**:

```bash
go install github.com/igolaizola/askimg/cmd/askimg@latest
```

Or you can download the binary from the [releases](https://github.com/igolaizola/askimg/releases)

## 📋 Requirements

You need to have a [Replicate](https://replicate.com/) account and a valid [API token](https://replicate.com/account/api-tokens).

## 🕹️ Usage

Using a configuration file in YAML format:

```bash
askimg --config askimg.conf
```

```yaml
# askimg.conf
token: REPLICATE_TOKEN
image: http://example.com/car.jpg
question: What is the color of the car?
```

Using environment variables (`ASKIMG` prefix, uppercase and underscores):

```bash
export ASKIMG_TOKEN=REPLICATE_TOKEN
export ASKIMG_IMAGE="http://example.com/car.jpg"
export ASKIMG_QUESTION="What is the color of the car?"
```

Using command line arguments:

```bash
askimg --token REPLICATE_TOKEN --image http://example.com/car.jpg --question "What is the color of the car?"
```

## 💖 Support

If you have found my code helpful, please give the repository a star ⭐

Additionally, if you would like to support my late-night coding efforts and the coffee that keeps me going, I would greatly appreciate a donation.

You can invite me for a coffee at ko-fi (0% fees):

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/igolaizola)

Or at buymeacoffee:

[![buymeacoffee](https://user-images.githubusercontent.com/11333576/223217083-123c2c53-6ab8-4ea8-a2c8-c6cb5d08e8d2.png)](https://buymeacoffee.com/igolaizola)

Donate to my PayPal:

[paypal.me/igolaizola](https://www.paypal.me/igolaizola)

Sponsor me on GitHub:

[github.com/sponsors/igolaizola](https://github.com/sponsors/igolaizola)

Or donate to any of my crypto addresses:

 - BTC `bc1qvuyrqwhml65adlu0j6l59mpfeez8ahdmm6t3ge`
 - ETH `0x960a7a9cdba245c106F729170693C0BaE8b2fdeD`
 - USDT (TRC20) `TD35PTZhsvWmR5gB12cVLtJwZtTv1nroDU`
 - USDC (BEP20) / BUSD (BEP20) `0x960a7a9cdba245c106F729170693C0BaE8b2fdeD`
 - Monero `41yc4R9d9iZMePe47VbfameDWASYrVcjoZJhJHFaK7DM3F2F41HmcygCrnLptS4hkiJARCwQcWbkW9k1z1xQtGSCAu3A7V4`

Thanks for your support!
