# 9. Metadata to identify unset bool values

Date: 2024-12-30

## Status

Accepted

## Context

Let's assume the boolean configuration item "force" is set by (the lower-priority) default configuration to "false".
Let's further assume that in the (higher-priority) user configuration, the value is not set at all.

The question is: What is the value of "force" in this case? We need to merge four layers of configuration
(as described in ADR-0008-configuration-layering.md), and need a method to handle unspecified boolean values.

The _old_ approach was the following (error-prone) one:

```golang
lowerPrioConf.Verbose = higherPrioConf.Verbose
```


## Decision

#### 1. Introduce metadata in the configuration struct to identify unset bool values.

```golang
type Config struct {
   OptionA bool
   OptionB bool
   OptionC bool
   // Metadata to track which fields were explicitly set
   SetFields map[string]bool
}
```

#### 2. Track Explicitly Set Fields
When populating the highConfig, explicitly mark fields as "set" in the SetFields map.

```golang
func NewConfig() Config {
   return Config{
      SetFields: make(map[string]bool),
   }
}

   highPrioConf.OptionA = false
   highPrioConf.SetFields["OptionA"] = true
```

#### 3. Merge Configurations
Use the SetFields map during merging to determine whether to use the value from highConfig or lowConfig.

```golang
func MergeConfigs(highConfig, lowConfig Config) Config {
   merged := Config{
      SetFields: make(map[string]bool),
   }

   if highConfig.SetFields["OptionA"] {
       merged.OptionA = highConfig.OptionA
   } else {
       merged.OptionA = lowConfig.OptionA
  }
```

#### 4. Enhancement: Utility Functions
We add a helper function to abstract common checking if a field was set.

```go
func (c Config) IsSet(field string) bool {
   return c.SetFields[field]
}
```

## Consequences

We can now safely merge configurations, even if certain higher-priority configurations (like files or flags) leave certain fields unset.

Downside is the slightly higher overhead when operating on configuration variables.