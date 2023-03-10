-- +migrate Up
--
-- name: SeedPosts :exec
--
INSERT INTO posts (
ID,
post_title,
post_slug,
post_content,
post_author,
post_image
)
VALUES (
        1,
        'Auto Draft',
        'auto-draft',
        'To be, or not to be: that is the question (Hamlet, Act 3, Scene 1). Romeo, Romeo! Wherefore art thou Romeo? (Romeo and Juliet, Act 2, Scene 2). Now is the winter of our discontent (Richard III, Act 1, Scene 1)',
        1,
        1
    ),
    (
        2,
        'Auto Draft 1',
        'auto-draft 1',
        'To be, or not to be: that is the question (Hamlet, Act 3, Scene 1). Romeo, Romeo! Wherefore art thou Romeo? (Romeo and Juliet, Act 2, Scene 2). Now is the winter of our discontent (Richard III, Act 1, Scene 1)',
        1,
        1
    );
--
-- name: SeedCategories :exec
--
INSERT INTO categories (
`ID`,
`category_title`,
`category_slug`,
`description`
)
VALUES (1, 'All categories', 'all-categories', '');
--
-- name: SeedUsers :exec
--
INSERT INTO users (
ID,
user_name,
user_email,
user_pass,
user_role,
display_name,
user_image
)
VALUES (
        1,
        'admin',
        'myemailaddress@emaildomain.com',
        '827ccb0eea8a706c4c34a16891f84e7b',
        '4',
        'admin',
        '1'
    );
--
-- name: SeedImages :exec
--
INSERT INTO images (
ID,
image_title,
image_url,
thumbnail_url,
height,
width
)
VALUES (
        1,
        'cat',
        'https://images.pexels.com/photos/7524926/pexels-photo-7524926.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750',
        'https://images.pexels.com/photos/7524926/pexels-photo-7524926.jpeg?auto=compress&cs=tinysrgb&w=400',
        '1260',
        '750'
    );
-- +migrate Down
TRUNCATE TABLE `posts `;
TRUNCATE TABLE `categories`;
TRUNCATE TABLE `users `;
TRUNCATE TABLE `images`;