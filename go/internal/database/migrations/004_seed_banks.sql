INSERT INTO banks (name, code, type, status) VALUES
    ('Bank Rakyat Indonesia (BRI)',     '002',   'bank',    'Active'),
    ('Bank Mandiri',                    '008',   'bank',    'Active'),
    ('Bank Negara Indonesia (BNI)',     '009',   'bank',    'Active'),
    ('Bank Tabungan Negara (BTN)',      '200',   'bank',    'Active'),
    ('Bank Central Asia (BCA)',         '014',   'bank',    'Active'),
    ('Bank Permata',                    '013',   'bank',    'Active'),
    ('Bank CIMB Niaga',                 '022',   'bank',    'Active'),
    ('Gopay',                           'gopay', 'wallet',  'Active'),
    ('OVO',                             'ovo',   'wallet',  'Active'),
    ('DANA',                            'dana',  'wallet',  'Active')
ON CONFLICT (code) DO NOTHING;