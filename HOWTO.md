ssllabs-scan -hostfile erc-list | jq -j -c . >erc-list-raw.json


v = `SITE-raw.json`
site_id = SELECT id FROM sites WHERE name = v.host
INSERT INTO runs ( site_id, d_when, data) VALUES ( N, now(), v::text)
