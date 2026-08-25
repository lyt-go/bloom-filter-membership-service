package filtercache

type Draft struct{ ID, Strategy string }
type Cache struct{ items map[string]*Draft }

func New() *Cache                             { return &Cache{items: map[string]*Draft{}} }
func (c *Cache) Put(d *Draft)                 { c.items[d.ID] = d }
func (c *Cache) Get(id string) (*Draft, bool) { d, ok := c.items[id]; return d, ok }
