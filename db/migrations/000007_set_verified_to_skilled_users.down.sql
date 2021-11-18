UPDATE users
SET verified_reviewer = false
WHERE rank = 'Immortal'
OR rank = 'Divine 5';

