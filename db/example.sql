INSERT INTO blogs (id, title, content, author, posted_at) VALUES
  (1, 'Hello World!', 'This is an example text.', 1, '2000-01-01 00:00'),
  (2, 'Hello more world!', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec lorem tellus, tincidunt sed pellentesque sit amet, sollicitudin sed nisl. Suspendisse potenti. Integer convallis ex vel libero molestie maximus. Sed efficitur et purus ut tincidunt. Curabitur mattis rutrum dolor a vehicula. Ut enim justo, lacinia sit amet pharetra mollis, congue quis ligula. Curabitur in nisl varius leo consequat imperdiet. In hac habitasse platea dictumst. Suspendisse imperdiet eget mi at lacinia. In eu felis quis ligula rutrum auctor ac in dui. Donec sit amet lorem a neque viverra dignissim quis aliquam libero. Mauris sed libero ligula. Cras varius vestibulum dui, ut finibus lectus lacinia ut. Sed non arcu ut ligula iaculis congue et sit amet nulla. Nunc et purus ex. Nulla a justo lacus.
Ut malesuada turpis massa, id ornare erat luctus ut. Donec risus nibh, fringilla vitae orci quis, egestas tristique risus. Nulla pharetra leo eros, ultricies ultrices erat eleifend non. Curabitur lacinia quis orci sed feugiat. Morbi sit amet pulvinar sem, vel ultrices elit. Sed a massa et dui bibendum convallis. Nam finibus dui felis, vel ullamcorper libero suscipit nec. Morbi sem nulla, pulvinar sed ultrices eget, laoreet at dui. Vivamus scelerisque egestas auctor. Praesent ultricies sed nisl ut cursus. Proin ut rhoncus lacus. Sed sem quam, venenatis quis cursus non, ornare ut ligula. Suspendisse vestibulum ut dolor vitae sodales. Aliquam vitae tincidunt dolor, eget semper tortor. Pellentesque id auctor nisi.
Nam imperdiet porttitor consectetur. Duis sed ornare diam, id maximus arcu. Sed vehicula arcu ut massa egestas porttitor. Nullam eget sapien malesuada, finibus metus vitae, egestas nisl. Nulla mattis dui rutrum risus varius tincidunt. Nulla commodo finibus ligula, eget mattis lacus blandit vitae. Pellentesque eget commodo justo, quis maximus ex. Aliquam rutrum arcu in gravida volutpat. Nulla imperdiet lacus eget sapien eleifend fermentum. Suspendisse potenti. Aliquam ligula ante, posuere ut ornare at, maximus id sapien. Nunc tristique quam diam, vitae consequat elit consectetur in. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Praesent rutrum vulputate felis, tincidunt posuere mauris faucibus dignissim. Pellentesque in diam sit amet neque varius faucibus.
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas sed malesuada quam. Aenean semper a odio id malesuada. Duis magna tortor, rhoncus non dui at, luctus varius diam. Praesent pellentesque odio diam, quis faucibus turpis blandit quis. Suspendisse erat quam, suscipit eget accumsan ac, commodo eget sem. Donec gravida urna nibh, sed varius est ultricies sit amet. Morbi gravida et erat non accumsan. Nulla placerat, arcu sed convallis tincidunt, arcu dui ullamcorper erat, ac pulvinar leo arcu vitae mi. Aliquam volutpat est nulla, id commodo turpis dapibus ac. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed posuere viverra arcu, sit amet tristique purus tincidunt in. Aenean lorem turpis, imperdiet at luctus eget, gravida in nibh.
Mauris ac interdum orci. Quisque sed metus libero. Donec at venenatis risus, a tempor lorem. Nulla quis eros urna. Fusce eu lectus nunc. Integer vitae tempor metus, vitae feugiat ex. Suspendisse a libero ultricies, vestibulum neque et, facilisis dui. Mauris sit amet neque dapibus, hendrerit urna at, vestibulum velit. Donec euismod iaculis ante, et scelerisque elit sollicitudin eget. Morbi nec dui ut urna elementum consectetur in sed nunc. Praesent volutpat neque fermentum metus accumsan iaculis. Nullam nec sodales felis, nec auctor purus. Vivamus dui magna, cursus in mauris sed, fringilla tempor risus.',
   1, '2000-01-02 00:00');

INSERT INTO artists (id, name, bio, added, added_by) VALUES
  (1, 'Led Zeppelin', 'Some American Band', now(), 1),
  (2, 'deadmau5', 'Some Canadian producer of electronic music', now(), 1),
  (3, 'No Mana', 'Some electronic dude', now(), 1);

INSERT INTO artist_tags (id, tag) VALUES
  (1, 'rock'),
  (2, '70s'),
  (3, '80s'),
  (4, '2000s'),
  (5, '2010s'),
  (6, 'usa'),
  (7, 'canada'),
  (8, 'edm'),
  (9, 'techno');

INSERT INTO artist_tags_artists (artist, tag) VALUES
  (1, 1),
  (1, 2),
  (1, 3),
  (1, 6),
  (2, 4),
  (2, 5),
  (2, 7),
  (2, 8),
  (2, 9),
  (3, 5),
  (3, 8),
  (3, 9);

INSERT INTO artist_aliases (artist, alias, added, added_by)
VALUES (2, 'testpilot', now(), 1);

INSERT INTO release_group_tags (id, tag) VALUES
  (1, 'edm'),
  (2, 'techno'),
  (3, 'electronic');

INSERT INTO release_groups (id, name, type, release_date, added, added_by)
VALUES (1, '4x4=12', 0, '2010-12-03', now(), 1);

INSERT INTO release_group_tags_release_groups (release_group, tag) VALUES
  (1, 1),
  (1, 2),
  (1, 3);

INSERT INTO release_groups_artists (release_group, artist, role)
VALUES (1, 2, 0);