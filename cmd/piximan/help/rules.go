package help

import "fmt"

// TODO: make these specialized manuals for infer id and path patterns and such

const rulesHelp = //
`Rules are used to filter wich works should be downloaded and defined in YAML format.
All rules are optional, and if multiple rules are defined, the work should match all of
them to be downloaded (AND). Any array matches any of its elements (OR).

If you download bookmarks with --low-meta flag, be aware that rules, related to the
missing metadata fields will be ignored. You need to download full metadata to match them.

Here's an example of all available rules:

  ids:                       [12345, 23456]
  not_ids:                   [34567, 45678]
  title_contains:            ['cute', 'cat']
  title_not_contains:        ['ugly', 'dog']
  title_regexp:              '^.*[0-9]+$'
  kinds:                     ['illust', 'manga', 'ugoira', 'novel']
  description_contains:      ['hello', 'world']
  description_not_contains:  ['goodbye', 'universe']
  description_regexp:        '^.*[0-9]+$'
  user_ids:                  [12345, 23456]
  not_user_ids:              [34567, 45678]
  user_names:                ['fekoneko', 'somecoolartist']
  not_user_names:            ['notsocoolartist', 'notme']
  restrictions:              ['none', 'R-18', 'R-18G']
  ai:                        false
  original:                  true
  pages_less_than:           50
  pages_more_than:           3
  views_less_than:           10000
  views_more_than:           1000
  bookmarks_less_than:       1000
  bookmarks_more_than:       100
  likes_less_than:           500
  likes_more_than:           50
  comments_less_than:        10
  comments_more_than:        2
  uploaded_before:           2022-01-01T00:00:00Z00:00
  uploaded_after:            2010-01-01T00:00:00Z00:00
  series:                    true
  series_ids:                [12345, 23456]
  not_series_ids:            [34567, 45678]
  series_title_contains:     ['cute', 'cat']
  series_title_not_contains: ['ugly', 'dog']
  series_title_regexp:       '^.*[0-9]+$'
  tags:                      ['お気に入り', '東方']
  not_tags:                  ['おっぱい', 'AI生成']
`

func RunRules() {
	fmt.Print(rulesHelp)
}
