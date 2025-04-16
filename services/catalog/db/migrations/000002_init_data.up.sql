-- @formatter:off
INSERT INTO publisher (id, name, address) VALUES (1, 'O''Reilly Media', '1005 Gravenstein Highway North Sebastopol, CA 95472, USA');
INSERT INTO publisher (id, name, address) VALUES (2, 'Manning Publications', '20 Baldwin Road PO Box 261 Shelter Island, NY 11964, USA');
INSERT INTO publisher (id, name, address) VALUES (3, 'Apress Media', 'One New York Plaza, Suite 4600 New York, NY 10004-1562, USA');
INSERT INTO publisher (id, name, address) VALUES (4, 'Addison-Wesley Professional', NULL);
INSERT INTO publisher (id, name, address) VALUES (5, 'Packt Publishing', 'Grosvenor House, 11 St Paul''s Square, Birmingham, B3 1RB, UK');
INSERT INTO publisher (id, name, address) VALUES (6, 'Pearson Publishing', NULL);

INSERT INTO book (id, title, category, isbn, price, publisher_id, description)
VALUES (1, 'Clean Code: A Handbook of Agile Software Craftsmanship', 18, '9780132350884', 47.99, 6, 'Even bad code can function. But if code isn''t clean, it can bring a development organization to its knees. Every year, countless hours and significant resources are lost because of poorly written code. But it doesn''t have to be that way.');
INSERT INTO book (id, title, category, isbn, price, publisher_id, description)
VALUES (2, 'Clean Architecture: A Craftsman''s Guide to Software Structure and Design', 19, '9780134494326', 49.99, 6, 'By applying universal rules of software architecture, you can dramatically improve developer productivity throughout the life of any software system.');
INSERT INTO book (id, title, category, isbn, price, publisher_id, description)
VALUES (3, 'Designing Data-Intensive Applications', 20, '9781449373320', 59.99, 1, 'Data is at the center of many challenges in system design today. Difficult issues need to be figured out, such as scalability, consistency, reliability, efficiency, and maintainability.');
INSERT INTO book (id, title, category, isbn, price, publisher_id, description)
VALUES (4, 'The Pragmatic Programmer: your journey to mastery, 20th Anniversary Edition, 2nd Edition', 21, '9780135957059', 54.99, 4, 'The Pragmatic Programmer is one of those rare tech books you''ll read, re-read, and read again over the years. Whether you''re new to the field or an experienced practitioner, you''ll come away with fresh insights each and every time.');
INSERT INTO book (id, title, category, isbn, price, publisher_id, description)
VALUES (5, 'Refactoring: Improving the Design of Existing Code', 22 , '9780134757599', 65.99, 4, 'For more than twenty years, experienced programmers worldwide have relied on Martin Fowler’s Refactoring to improve the design of existing code and to enhance software maintainability, as well as to make existing code easier to understand.');
INSERT INTO book (id, title, category, isbn, price, publisher_id, description)
VALUES (6, 'Kubernetes in Action', 23, '9781617293726', 59.99, 2, 'Kubernetes in Action is a comprehensive guide to effectively developing and running applications in a Kubernetes environment. Before diving into Kubernetes, the book gives an overview of container technologies like Docker, including how to build containers, so that even readers who haven''t used these technologies before can get up and running.');

INSERT INTO author (id, name) VALUES (1, 'Robert C. Martin');
INSERT INTO author (id, name) VALUES (2, 'Martin Fowler');
INSERT INTO author (id, name) VALUES (3, 'Kent Beck');
INSERT INTO author (id, name) VALUES (4, 'Joshua Bloch');
INSERT INTO author (id, name) VALUES (5, 'Martin Kleppmann');
INSERT INTO author (id, name) VALUES (6, 'Adam Freeman');
INSERT INTO author (id, name) VALUES (7, 'Andrew Hunt');
INSERT INTO author (id, name) VALUES (8, 'Mark Michaelis');
INSERT INTO author (id, name) VALUES (9, 'Marko Luksa');
INSERT INTO author (id, name) VALUES (10, 'David Thomas');
INSERT INTO author (id, name) VALUES (11, 'Andrew Hunt');

INSERT INTO book_author_mapping (book_id, author_id) VALUES (1, 1);
INSERT INTO book_author_mapping (book_id, author_id) VALUES (2, 1);
INSERT INTO book_author_mapping (book_id, author_id) VALUES (3, 5);
INSERT INTO book_author_mapping (book_id, author_id) VALUES (4, 10);
INSERT INTO book_author_mapping (book_id, author_id) VALUES (4, 11);
INSERT INTO book_author_mapping (book_id, author_id) VALUES (5, 2);
INSERT INTO book_author_mapping (book_id, author_id) VALUES (6, 9);
