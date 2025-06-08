import itertools
import random
import time

import yaml  # type: ignore

# === Настройки ===
CONFIG_FILE = "tools/metrics_config.yaml"
OUTPUT_FILE = "tools/metrics.txt"
VICTORIA_URL = "http://localhost:8428/api/v1/import/prometheus"

def load_config(path):
    with open(path) as f:
        return yaml.safe_load(f)

def generate_label_combinations(label_dict):
    if not label_dict:
        yield {}
        return
    keys = label_dict.keys()
    values = label_dict.values()
    for combo in itertools.product(*values):
        yield dict(zip(keys, combo))

def format_metric(name, labels, value, timestamp):
    label_str = ",".join(f'{k}="{v}"' for k, v in labels.items())
    return f'{name}{{{label_str}}} {value:.3f} {timestamp}'

def generate_metric_lines(metric):
    name = metric["name"]
    mtype = metric.get("type", "gauge")
    label_defs = metric.get("labels", {})
    min_val, max_val = metric["value_range"]
    interval = metric.get("interval", 60)
    points = metric.get("points", 1)

    lines = []
    for labels in generate_label_combinations(label_defs):
        current = random.uniform(min_val, max_val) if mtype == "counter" else 0.0
        for i in range(points):
            ts = int(time.time()) - (points - i) * interval
            if mtype == "counter":
                current += random.uniform(min_val, max_val)
            else:  # gauge
                current = random.uniform(min_val, max_val)
            line = format_metric(name, labels, current, ts)
            lines.append(line)
    return lines

def main():
    config = load_config(CONFIG_FILE)
    all_lines = []
    for metric in config["metrics"]:
        all_lines.extend(generate_metric_lines(metric))

    with open(OUTPUT_FILE, "w") as f:
        f.write("\n".join(all_lines) + "\n")

    print(f"✅ Generated {len(all_lines)} metrics into {OUTPUT_FILE}")

    # with open(OUTPUT_FILE, "rb") as f:
    #     resp = requests.post(VICTORIA_URL, data=f)
    #     print(f"📤 Push status: {resp.status_code} {resp.reason}")

if __name__ == "__main__":
    main()
