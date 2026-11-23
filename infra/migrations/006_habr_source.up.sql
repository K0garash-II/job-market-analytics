INSERT INTO sources (name, url) VALUES ('habr', 'https://career.habr.com') ON CONFLICT (name) DO NOTHING;
