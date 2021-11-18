UPDATE users
SET verified_reviewer = true
WHERE rank = 'Immortal'
OR rank = 'Divine 5';
