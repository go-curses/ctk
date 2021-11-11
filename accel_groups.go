package ctk

//
// type AccelGroups interface {
// 	Activate(object Object, accelKey cdk.Key, accelMods ModifierType) (value bool)
// 	FromObject(object Object) (value []*AccelGroup)
// }
//
// type CAccelGroups []*CAccelGroup
//
// func NewAccelGroups() *CAccelGroups {
// 	return &CAccelGroups{}
// }
//
// func (a *CAccelGroups) Activate(object Object, accelKey cdk.Key, accelMods ModifierType) (value bool) {
// 	for _, ag := range *a {
// 		entries := ag.Query(accelKey, accelMods)
// 		if len(entries) > 0 {
// 			entries[0].Closure()
// 		}
// 	}
// 	return
// }
