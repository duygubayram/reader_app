-- =========================
-- Books table
-- =========================

CREATE TABLE IF NOT EXISTS books (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    author TEXT NOT NULL,
    year INTEGER,
    language TEXT,
    publisher TEXT,
    total_pages INTEGER
);

-- =========================
-- Fake book data (50 books) - chatgpt generated
-- =========================

INSERT OR IGNORE INTO books (id, name, author, year, publisher, language, total_pages) VALUES
(1, 'The Silent Horizon', 'Ava Mitchell', 2012, 'Northwind Press', 'English', 384),
(2, 'Ashes of Tomorrow', 'Liam Carter', 2015, 'Red Maple Books', 'English', 421),
(3, 'The Clockmaker’s Paradox', 'Eleanor Finch', 2018, 'Ironleaf Publishing', 'English', 356),
(4, 'Beyond the Riverbend', 'Noah Whitaker', 2010, 'Stonebridge House', 'English', 298),
(5, 'Fragments of Light', 'Isabella Moreau', 2016, 'Étoile Éditions', 'French', 332),
(6, 'The Last Archivist', 'Marcus Hale', 2020, 'Blackwell & Co.', 'English', 467),
(7, 'Paper Cities', 'Sofia Lindström', 2014, 'Nordic Ink', 'Swedish', 289),
(8, 'A Theory of Forgotten Things', 'Julian Rowe', 2019, 'Helix Press', 'English', 410),
(9, 'Winter Over Caldera', 'Hannah Brooks', 2011, 'Summit Lane', 'English', 276),
(10, 'The Glass Orchard', 'Theo Alvarez', 2017, 'Sunstone Publishing', 'Spanish', 345),

(11, 'Maps for the Unlost', 'Priya Nandakumar', 2021, 'Lotus River Press', 'English', 392),
(12, 'Ink Beneath the Skin', 'Ronan Price', 2013, 'Cinder House', 'English', 318),
(13, 'When Statues Dream', 'Alessandro Ricci', 2009, 'Via Roma Books', 'Italian', 264),
(14, 'The Narrow Season', 'Emily Zhao', 2018, 'Paper Crane Publishing', 'English', 351),
(15, 'Letters Never Sent', 'Marta Kowalska', 2012, 'Baltic Words', 'Polish', 287),
(16, 'The Echo Cartographer', 'Samuel Keene', 2022, 'Wayfinder Press', 'English', 455),
(17, 'Dust and Other Silences', 'Omar El-Tayeb', 2016, 'Desert Palm Books', 'Arabic', 309),
(18, 'Anatomy of a Firefly', 'Clara Voss', 2019, 'Moonwell Press', 'English', 334),
(19, 'The Long Province', 'Daniel Okoye', 2011, 'Crosswind Publishing', 'English', 401),
(20, 'Blue Hours in Kyoto', 'Rei Nakamura', 2015, 'Hoshino Books', 'Japanese', 258),

(21, 'The Sound of Distant Bells', 'Margaret Hill', 2008, 'Elder Grove', 'English', 372),
(22, 'After the Cedar Falls', 'Jonah Peterson', 2014, 'Timberline Press', 'English', 319),
(23, 'Saltwater Arithmetic', 'Nina Calder', 2020, 'Driftwood Editions', 'English', 386),
(24, 'The Unfinished Atlas', 'Victor Laurent', 2017, 'Meridian House', 'French', 428),
(25, 'A Study in Hollow Places', 'Beatrice Young', 2021, 'Ravencrest Press', 'English', 399),
(26, 'How Mountains Remember', 'Lucas Fernández', 2013, 'Andes Ink', 'Spanish', 344),
(27, 'The Color of Returning', 'Yara Haddad', 2019, 'Olive Branch Books', 'Arabic', 361),
(28, 'Small Gods of Concrete', 'Patrick Doyle', 2016, 'Iron Alley Press', 'English', 295),
(29, 'The Third Weather', 'Helena Novak', 2018, 'Silver Birch Publishing', 'Czech', 327),
(30, 'Notes from a Vanishing Shore', 'Caleb Morris', 2012, 'Low Tide Press', 'English', 282),

(31, 'The City That Waited', 'Irene Park', 2022, 'Neon Harbor', 'English', 447),
(32, 'What the Moss Keeps', 'Frederik Olsen', 2010, 'Green Fjord Press', 'Danish', 271),
(33, 'The Second Library of Babel', 'Arthur Klein', 2019, 'Parallax Books', 'English', 503),
(34, 'A Brief History of Falling', 'Lydia Chen', 2015, 'Skybound Press', 'English', 318),
(35, 'Wind, Stone, Threshold', 'Mikhail Petrov', 2007, 'Volga House', 'Russian', 389),
(36, 'The Shape of Abandoned Roads', 'Oliver Grant', 2018, 'Farway Publishing', 'English', 362),
(37, 'The Sea Is Not a Mirror', 'Lucía Morales', 2020, 'Azul Mar Editions', 'Spanish', 341),
(38, 'Rooms Without North', 'Anika Bose', 2016, 'Eastline Press', 'English', 304),
(39, 'The Persistence of Smoke', 'Thomas Reed', 2011, 'Ashfall Publishing', 'English', 376),
(40, 'Beneath Quiet Signals', 'Jonas Weber', 2014, 'Rhinegold Books', 'German', 333),

(41, 'The Last Season of Maps', 'Naomi Feldman', 2021, 'Compass Rose Press', 'English', 412),
(42, 'If Stones Could Speak', 'Eamon Walsh', 2009, 'Cloverfield House', 'English', 291),
(43, 'A Manual for Temporary Lives', 'Sanjay Rao', 2017, 'Open Palm Press', 'English', 359),
(44, 'Night Letters to the Coast', 'Phoebe Lang', 2013, 'Lighthouse Ink', 'English', 274),
(45, 'The Orchard After Snow', 'Katerina Ivanova', 2018, 'White Ember Press', 'Russian', 336),
(46, 'Distances We Inherit', 'Michael Osei', 2020, 'New Meridian Press', 'English', 398),
(47, 'The Grammar of Tides', 'Elena Rossi', 2016, 'Blue Current Publishing', 'Italian', 321),
(48, 'Cities Built of Breath', 'Farah Suleiman', 2022, 'Minaret Books', 'Arabic', 452),
(49, 'An Inventory of Absences', 'George Whitman', 2011, 'Fieldnote Press', 'English', 365),
(50, 'The Year We Learned to Wait', 'Amara Singh', 2019, 'Stillwater Press', 'English', 387);
