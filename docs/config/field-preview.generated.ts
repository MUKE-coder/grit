// Code generated from internal/generate. DO NOT EDIT.
//
// What each field type produces across the stack, asked of the same methods
// the generator uses: Field.GoType, GORMTag, TSType, ZodType and
// FormFieldType. Regenerate with:
//
//	go test ./internal/generate -run TestFieldPreviewDataIsCurrent -update
//
// TestFieldPreviewDataIsCurrent fails when this file drifts from them, so the
// field builder cannot describe a framework that no longer exists.
//
// required and unique each get their own shape, because they change the
// answer: Required adds .min(1, "Required") to a string and drops the
// .optional(), and Unique adds uniqueIndex to the column. They are left out
// when they are the same as base.

export type FieldShape = {
  /** The struct field's type, as the model declares it. */
  go: string
  /** The struct tag that decides the column, empty when it needs none. */
  gorm: string
  /** The TypeScript type in packages/shared. */
  ts: string
  /** The rule the create and update schemas carry. */
  zod: string
  /** The admin form control, empty when the server fills the column. */
  control: string
}

export type FieldPreview = {
  /** The name written in --fields, such as "money". */
  type: string
  /** With neither modifier. */
  base: FieldShape
  /** Absent when the modifier changes nothing. */
  required?: FieldShape
  unique?: FieldShape
  requiredUnique?: FieldShape
  /** The one thing the shape does not say. */
  note?: string
}

/** The shape for a field with these modifiers, falling back to base. */
export function shapeFor(
  preview: FieldPreview,
  modifiers: { required?: boolean; unique?: boolean },
): FieldShape {
  const { required, unique } = modifiers
  if (required && unique) return preview.requiredUnique ?? preview.base
  if (required) return preview.required ?? preview.base
  if (unique) return preview.unique ?? preview.base
  return preview.base
}

export const FIELD_PREVIEWS: Record<string, FieldPreview> = {
  belongs_to: {
    type: 'belongs_to',
    base: {
      go: 'string',
      gorm: 'size:36;index',
      ts: 'string',
      zod: 'z.uuid("Invalid ID")',
      control: 'relationship-select',
    },
    unique: {
      go: 'string',
      gorm: 'size:36;index;uniqueIndex',
      ts: 'string',
      zod: 'z.uuid("Invalid ID")',
      control: 'relationship-select',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'size:36;index;uniqueIndex',
      ts: 'string',
      zod: 'z.uuid("Invalid ID")',
      control: 'relationship-select',
    },
    note: 'Writes the foreign key column and the association, and the admin gets a picker.',
  },
  bool: {
    type: 'bool',
    base: {
      go: 'bool',
      gorm: '',
      ts: 'boolean',
      zod: 'z.boolean().optional()',
      control: 'toggle',
    },
    required: {
      go: 'bool',
      gorm: '',
      ts: 'boolean',
      zod: 'z.boolean()',
      control: 'toggle',
    },
    unique: {
      go: 'bool',
      gorm: 'uniqueIndex',
      ts: 'boolean',
      zod: 'z.boolean().optional()',
      control: 'toggle',
    },
    requiredUnique: {
      go: 'bool',
      gorm: 'uniqueIndex',
      ts: 'boolean',
      zod: 'z.boolean()',
      control: 'toggle',
    },
  },
  check: {
    type: 'check',
    base: {
      go: 'datatypes.JSONSlice[string]',
      gorm: 'type:json',
      ts: '("draft" | "sent")[]',
      zod: 'z.array(z.enum(["draft", "sent"])).optional()',
      control: 'checkbox-group',
    },
    unique: {
      go: 'datatypes.JSONSlice[string]',
      gorm: 'type:json;uniqueIndex',
      ts: '("draft" | "sent")[]',
      zod: 'z.array(z.enum(["draft", "sent"])).optional()',
      control: 'checkbox-group',
    },
    requiredUnique: {
      go: 'datatypes.JSONSlice[string]',
      gorm: 'type:json;uniqueIndex',
      ts: '("draft" | "sent")[]',
      zod: 'z.array(z.enum(["draft", "sent"])).optional()',
      control: 'checkbox-group',
    },
    note: 'Several choices in one JSON column, not a join table. The union and the enum above show an example; yours take their place.',
  },
  color: {
    type: 'color',
    base: {
      go: 'string',
      gorm: 'size:7',
      ts: 'string',
      zod: 'ColorSchema.or(z.literal("")).optional()',
      control: 'color',
    },
    required: {
      go: 'string',
      gorm: 'size:7',
      ts: 'string',
      zod: 'ColorSchema',
      control: 'color',
    },
    unique: {
      go: 'string',
      gorm: 'size:7;uniqueIndex',
      ts: 'string',
      zod: 'ColorSchema.or(z.literal("")).optional()',
      control: 'color',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'size:7;uniqueIndex',
      ts: 'string',
      zod: 'ColorSchema',
      control: 'color',
    },
  },
  country: {
    type: 'country',
    base: {
      go: 'string',
      gorm: 'size:2',
      ts: 'string',
      zod: 'CountrySchema.or(z.literal("")).optional()',
      control: 'country',
    },
    required: {
      go: 'string',
      gorm: 'size:2',
      ts: 'string',
      zod: 'CountrySchema',
      control: 'country',
    },
    unique: {
      go: 'string',
      gorm: 'size:2;uniqueIndex',
      ts: 'string',
      zod: 'CountrySchema.or(z.literal("")).optional()',
      control: 'country',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'size:2;uniqueIndex',
      ts: 'string',
      zod: 'CountrySchema',
      control: 'country',
    },
  },
  date: {
    type: 'date',
    base: {
      go: '*jsontime.Date',
      gorm: 'type:date',
      ts: 'string | null',
      zod: 'z.string().nullable().optional()',
      control: 'date',
    },
    required: {
      go: '*jsontime.Date',
      gorm: 'type:date',
      ts: 'string | null',
      zod: 'z.string().nullable()',
      control: 'date',
    },
    unique: {
      go: '*jsontime.Date',
      gorm: 'type:date;uniqueIndex',
      ts: 'string | null',
      zod: 'z.string().nullable().optional()',
      control: 'date',
    },
    requiredUnique: {
      go: '*jsontime.Date',
      gorm: 'type:date;uniqueIndex',
      ts: 'string | null',
      zod: 'z.string().nullable()',
      control: 'date',
    },
    note: 'A day with no time, so it does not move when a reader is in another timezone.',
  },
  datetime: {
    type: 'datetime',
    base: {
      go: '*jsontime.DateTime',
      gorm: '',
      ts: 'string | null',
      zod: 'z.string().nullable().optional()',
      control: 'datetime',
    },
    required: {
      go: '*jsontime.DateTime',
      gorm: '',
      ts: 'string | null',
      zod: 'z.string().nullable()',
      control: 'datetime',
    },
    unique: {
      go: '*jsontime.DateTime',
      gorm: 'uniqueIndex',
      ts: 'string | null',
      zod: 'z.string().nullable().optional()',
      control: 'datetime',
    },
    requiredUnique: {
      go: '*jsontime.DateTime',
      gorm: 'uniqueIndex',
      ts: 'string | null',
      zod: 'z.string().nullable()',
      control: 'datetime',
    },
  },
  domain: {
    type: 'domain',
    base: {
      go: 'string',
      gorm: 'size:253',
      ts: 'string',
      zod: 'DomainSchema.or(z.literal("")).optional()',
      control: 'domain',
    },
    required: {
      go: 'string',
      gorm: 'size:253',
      ts: 'string',
      zod: 'DomainSchema',
      control: 'domain',
    },
    unique: {
      go: 'string',
      gorm: 'size:253;uniqueIndex',
      ts: 'string',
      zod: 'DomainSchema.or(z.literal("")).optional()',
      control: 'domain',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'size:253;uniqueIndex',
      ts: 'string',
      zod: 'DomainSchema',
      control: 'domain',
    },
  },
  email: {
    type: 'email',
    base: {
      go: 'string',
      gorm: 'size:254',
      ts: 'string',
      zod: 'EmailSchema.or(z.literal("")).optional()',
      control: 'email',
    },
    required: {
      go: 'string',
      gorm: 'size:254',
      ts: 'string',
      zod: 'EmailSchema',
      control: 'email',
    },
    unique: {
      go: 'string',
      gorm: 'size:254;uniqueIndex',
      ts: 'string',
      zod: 'EmailSchema.or(z.literal("")).optional()',
      control: 'email',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'size:254;uniqueIndex',
      ts: 'string',
      zod: 'EmailSchema',
      control: 'email',
    },
  },
  file: {
    type: 'file',
    base: {
      go: '*files.FileRef',
      gorm: 'type:json',
      ts: 'FileRef | null',
      zod: 'FileRefSchema.nullable().optional()',
      control: 'file',
    },
    required: {
      go: '*files.FileRef',
      gorm: 'type:json',
      ts: 'FileRef | null',
      zod: 'FileRefSchema',
      control: 'file',
    },
    unique: {
      go: '*files.FileRef',
      gorm: 'type:json;uniqueIndex',
      ts: 'FileRef | null',
      zod: 'FileRefSchema.nullable().optional()',
      control: 'file',
    },
    requiredUnique: {
      go: '*files.FileRef',
      gorm: 'type:json;uniqueIndex',
      ts: 'FileRef | null',
      zod: 'FileRefSchema',
      control: 'file',
    },
    note: 'Stores a reference, not the bytes. Deleting the row releases it to the orphan sweep.',
  },
  files: {
    type: 'files',
    base: {
      go: 'files.FileRefs',
      gorm: 'type:json',
      ts: 'FileRef[]',
      zod: 'z.array(FileRefSchema).default([])',
      control: 'files',
    },
    unique: {
      go: 'files.FileRefs',
      gorm: 'type:json;uniqueIndex',
      ts: 'FileRef[]',
      zod: 'z.array(FileRefSchema).default([])',
      control: 'files',
    },
    requiredUnique: {
      go: 'files.FileRefs',
      gorm: 'type:json;uniqueIndex',
      ts: 'FileRef[]',
      zod: 'z.array(FileRefSchema).default([])',
      control: 'files',
    },
    note: 'A list of references. The same release rules as file.',
  },
  float: {
    type: 'float',
    base: {
      go: 'float64',
      gorm: '',
      ts: 'number',
      zod: 'z.number().optional()',
      control: 'number',
    },
    required: {
      go: 'float64',
      gorm: '',
      ts: 'number',
      zod: 'z.number()',
      control: 'number',
    },
    unique: {
      go: 'float64',
      gorm: 'uniqueIndex',
      ts: 'number',
      zod: 'z.number().optional()',
      control: 'number',
    },
    requiredUnique: {
      go: 'float64',
      gorm: 'uniqueIndex',
      ts: 'number',
      zod: 'z.number()',
      control: 'number',
    },
  },
  int: {
    type: 'int',
    base: {
      go: 'int',
      gorm: '',
      ts: 'number',
      zod: 'z.number().int().optional()',
      control: 'number',
    },
    required: {
      go: 'int',
      gorm: '',
      ts: 'number',
      zod: 'z.number().int()',
      control: 'number',
    },
    unique: {
      go: 'int',
      gorm: 'uniqueIndex',
      ts: 'number',
      zod: 'z.number().int().optional()',
      control: 'number',
    },
    requiredUnique: {
      go: 'int',
      gorm: 'uniqueIndex',
      ts: 'number',
      zod: 'z.number().int()',
      control: 'number',
    },
  },
  json: {
    type: 'json',
    base: {
      go: 'datatypes.JSON',
      gorm: '',
      ts: 'unknown',
      zod: 'JsonValueSchema.optional()',
      control: 'json',
    },
    required: {
      go: 'datatypes.JSON',
      gorm: '',
      ts: 'unknown',
      zod: 'JsonValueSchema',
      control: 'json',
    },
    unique: {
      go: 'datatypes.JSON',
      gorm: 'uniqueIndex',
      ts: 'unknown',
      zod: 'JsonValueSchema.optional()',
      control: 'json',
    },
    requiredUnique: {
      go: 'datatypes.JSON',
      gorm: 'uniqueIndex',
      ts: 'unknown',
      zod: 'JsonValueSchema',
      control: 'json',
    },
    note: 'Stored as JSON and typed as unknown, so nothing downstream assumes a shape.',
  },
  many_to_many: {
    type: 'many_to_many',
    base: {
      go: '[]string',
      gorm: '',
      ts: 'string[]',
      zod: 'z.array(z.uuid()).optional()',
      control: 'multi-relationship-select',
    },
    note: 'Writes a join table. The admin gets a multi-select.',
  },
  money: {
    type: 'money',
    base: {
      go: 'money.Money',
      gorm: '',
      ts: 'Money',
      zod: 'MoneySchema.optional()',
      control: 'money',
    },
    required: {
      go: 'money.Money',
      gorm: '',
      ts: 'Money',
      zod: 'MoneySchema',
      control: 'money',
    },
    requiredUnique: {
      go: 'money.Money',
      gorm: '',
      ts: 'Money',
      zod: 'MoneySchema',
      control: 'money',
    },
    note: 'Two columns: the amount in minor units and the currency. Never a float.',
  },
  one_to_one: {
    type: 'one_to_one',
    base: {
      go: 'string',
      gorm: 'size:36;uniqueIndex',
      ts: 'string',
      zod: 'z.string().optional()',
      control: 'text',
    },
    required: {
      go: 'string',
      gorm: 'size:36;uniqueIndex',
      ts: 'string',
      zod: 'z.string()',
      control: 'text',
    },
    unique: {
      go: 'string',
      gorm: 'size:36;uniqueIndex;uniqueIndex',
      ts: 'string',
      zod: 'z.string().optional()',
      control: 'text',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'size:36;uniqueIndex;uniqueIndex',
      ts: 'string',
      zod: 'z.string()',
      control: 'text',
    },
  },
  percent: {
    type: 'percent',
    base: {
      go: 'float64',
      gorm: 'type:decimal(5,2)',
      ts: 'number',
      zod: 'PercentSchema.nullable().optional()',
      control: 'percent',
    },
    required: {
      go: 'float64',
      gorm: 'type:decimal(5,2)',
      ts: 'number',
      zod: 'PercentSchema',
      control: 'percent',
    },
    unique: {
      go: 'float64',
      gorm: 'type:decimal(5,2);uniqueIndex',
      ts: 'number',
      zod: 'PercentSchema.nullable().optional()',
      control: 'percent',
    },
    requiredUnique: {
      go: 'float64',
      gorm: 'type:decimal(5,2);uniqueIndex',
      ts: 'number',
      zod: 'PercentSchema',
      control: 'percent',
    },
  },
  radio: {
    type: 'radio',
    base: {
      go: 'string',
      gorm: 'size:255',
      ts: '"draft" | "sent"',
      zod: 'z.enum(["draft", "sent"]).optional()',
      control: 'radio',
    },
    required: {
      go: 'string',
      gorm: 'size:255',
      ts: '"draft" | "sent"',
      zod: 'z.enum(["draft", "sent"])',
      control: 'radio',
    },
    unique: {
      go: 'string',
      gorm: 'size:255;uniqueIndex',
      ts: '"draft" | "sent"',
      zod: 'z.enum(["draft", "sent"]).optional()',
      control: 'radio',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'size:255;uniqueIndex',
      ts: '"draft" | "sent"',
      zod: 'z.enum(["draft", "sent"])',
      control: 'radio',
    },
    note: 'The same column as select, with a different control. The union and the enum above show an example; your choices take their place.',
  },
  rating: {
    type: 'rating',
    base: {
      go: 'int',
      gorm: '',
      ts: 'number',
      zod: 'ratingSchema(5).nullable().optional()',
      control: 'rating',
    },
    required: {
      go: 'int',
      gorm: '',
      ts: 'number',
      zod: 'ratingSchema(5)',
      control: 'rating',
    },
    unique: {
      go: 'int',
      gorm: 'uniqueIndex',
      ts: 'number',
      zod: 'ratingSchema(5).nullable().optional()',
      control: 'rating',
    },
    requiredUnique: {
      go: 'int',
      gorm: 'uniqueIndex',
      ts: 'number',
      zod: 'ratingSchema(5)',
      control: 'rating',
    },
  },
  richtext: {
    type: 'richtext',
    base: {
      go: 'string',
      gorm: 'type:text',
      ts: 'string',
      zod: 'z.string().optional()',
      control: 'richtext',
    },
    required: {
      go: 'string',
      gorm: 'type:text',
      ts: 'string',
      zod: 'z.string()',
      control: 'richtext',
    },
    unique: {
      go: 'string',
      gorm: 'type:text;uniqueIndex',
      ts: 'string',
      zod: 'z.string().optional()',
      control: 'richtext',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'type:text;uniqueIndex',
      ts: 'string',
      zod: 'z.string()',
      control: 'richtext',
    },
    note: 'Sanitised on write, so what is stored is already safe to render.',
  },
  select: {
    type: 'select',
    base: {
      go: 'string',
      gorm: 'size:255',
      ts: '"draft" | "sent"',
      zod: 'z.enum(["draft", "sent"]).optional()',
      control: 'select',
    },
    required: {
      go: 'string',
      gorm: 'size:255',
      ts: '"draft" | "sent"',
      zod: 'z.enum(["draft", "sent"])',
      control: 'select',
    },
    unique: {
      go: 'string',
      gorm: 'size:255;uniqueIndex',
      ts: '"draft" | "sent"',
      zod: 'z.enum(["draft", "sent"]).optional()',
      control: 'select',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'size:255;uniqueIndex',
      ts: '"draft" | "sent"',
      zod: 'z.enum(["draft", "sent"])',
      control: 'select',
    },
    note: 'The choices are a database CHECK as well as a dropdown. The union and the enum above show an example; yours take their place.',
  },
  slug: {
    type: 'slug',
    base: {
      go: 'string',
      gorm: 'size:255;uniqueIndex',
      ts: 'string',
      zod: 'z.string()',
      control: '',
    },
    note: 'Filled by the server from another column, so it is not on the form.',
  },
  string: {
    type: 'string',
    base: {
      go: 'string',
      gorm: 'size:255',
      ts: 'string',
      zod: 'z.string().optional()',
      control: 'text',
    },
    required: {
      go: 'string',
      gorm: 'size:255',
      ts: 'string',
      zod: 'z.string().min(1, "Required")',
      control: 'text',
    },
    unique: {
      go: 'string',
      gorm: 'size:255;uniqueIndex',
      ts: 'string',
      zod: 'z.string().optional()',
      control: 'text',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'size:255;uniqueIndex',
      ts: 'string',
      zod: 'z.string().min(1, "Required")',
      control: 'text',
    },
  },
  string_array: {
    type: 'string_array',
    base: {
      go: 'datatypes.JSONSlice[string]',
      gorm: 'type:json',
      ts: 'string[]',
      zod: 'z.array(z.string()).optional()',
      control: 'images',
    },
    unique: {
      go: 'datatypes.JSONSlice[string]',
      gorm: 'type:json;uniqueIndex',
      ts: 'string[]',
      zod: 'z.array(z.string()).optional()',
      control: 'images',
    },
    requiredUnique: {
      go: 'datatypes.JSONSlice[string]',
      gorm: 'type:json;uniqueIndex',
      ts: 'string[]',
      zod: 'z.array(z.string()).optional()',
      control: 'images',
    },
    note: 'One column holding a list, not a join table.',
  },
  tel: {
    type: 'tel',
    base: {
      go: 'string',
      gorm: 'size:32',
      ts: 'string',
      zod: 'PhoneSchema.or(z.literal("")).optional()',
      control: 'tel',
    },
    required: {
      go: 'string',
      gorm: 'size:32',
      ts: 'string',
      zod: 'PhoneSchema',
      control: 'tel',
    },
    unique: {
      go: 'string',
      gorm: 'size:32;uniqueIndex',
      ts: 'string',
      zod: 'PhoneSchema.or(z.literal("")).optional()',
      control: 'tel',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'size:32;uniqueIndex',
      ts: 'string',
      zod: 'PhoneSchema',
      control: 'tel',
    },
  },
  text: {
    type: 'text',
    base: {
      go: 'string',
      gorm: 'type:text',
      ts: 'string',
      zod: 'z.string().optional()',
      control: 'textarea',
    },
    required: {
      go: 'string',
      gorm: 'type:text',
      ts: 'string',
      zod: 'z.string()',
      control: 'textarea',
    },
    unique: {
      go: 'string',
      gorm: 'type:text;uniqueIndex',
      ts: 'string',
      zod: 'z.string().optional()',
      control: 'textarea',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'type:text;uniqueIndex',
      ts: 'string',
      zod: 'z.string()',
      control: 'textarea',
    },
  },
  time: {
    type: 'time',
    base: {
      go: 'string',
      gorm: 'size:5',
      ts: 'string',
      zod: 'TimeSchema.or(z.literal("")).optional()',
      control: 'time',
    },
    required: {
      go: 'string',
      gorm: 'size:5',
      ts: 'string',
      zod: 'TimeSchema',
      control: 'time',
    },
    unique: {
      go: 'string',
      gorm: 'size:5;uniqueIndex',
      ts: 'string',
      zod: 'TimeSchema.or(z.literal("")).optional()',
      control: 'time',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'size:5;uniqueIndex',
      ts: 'string',
      zod: 'TimeSchema',
      control: 'time',
    },
  },
  toggle: {
    type: 'toggle',
    base: {
      go: 'bool',
      gorm: '',
      ts: 'boolean',
      zod: 'z.boolean().optional()',
      control: 'toggle',
    },
    required: {
      go: 'bool',
      gorm: '',
      ts: 'boolean',
      zod: 'z.boolean()',
      control: 'toggle',
    },
    unique: {
      go: 'bool',
      gorm: 'uniqueIndex',
      ts: 'boolean',
      zod: 'z.boolean().optional()',
      control: 'toggle',
    },
    requiredUnique: {
      go: 'bool',
      gorm: 'uniqueIndex',
      ts: 'boolean',
      zod: 'z.boolean()',
      control: 'toggle',
    },
  },
  uint: {
    type: 'uint',
    base: {
      go: 'uint',
      gorm: '',
      ts: 'number',
      zod: 'z.number().int().nonnegative().optional()',
      control: 'number',
    },
    required: {
      go: 'uint',
      gorm: '',
      ts: 'number',
      zod: 'z.number().int().nonnegative()',
      control: 'number',
    },
    unique: {
      go: 'uint',
      gorm: 'uniqueIndex',
      ts: 'number',
      zod: 'z.number().int().nonnegative().optional()',
      control: 'number',
    },
    requiredUnique: {
      go: 'uint',
      gorm: 'uniqueIndex',
      ts: 'number',
      zod: 'z.number().int().nonnegative()',
      control: 'number',
    },
  },
  url: {
    type: 'url',
    base: {
      go: 'string',
      gorm: 'size:2048',
      ts: 'string',
      zod: 'UrlSchema.or(z.literal("")).optional()',
      control: 'url',
    },
    required: {
      go: 'string',
      gorm: 'size:2048',
      ts: 'string',
      zod: 'UrlSchema',
      control: 'url',
    },
    unique: {
      go: 'string',
      gorm: 'size:2048;uniqueIndex',
      ts: 'string',
      zod: 'UrlSchema.or(z.literal("")).optional()',
      control: 'url',
    },
    requiredUnique: {
      go: 'string',
      gorm: 'size:2048;uniqueIndex',
      ts: 'string',
      zod: 'UrlSchema',
      control: 'url',
    },
  },
  vector: {
    type: 'vector',
    base: {
      go: 'vectors.Vector',
      gorm: '',
      ts: 'number[]',
      zod: 'z.array(z.number()).optional()',
      control: 'text',
    },
    required: {
      go: 'vectors.Vector',
      gorm: '',
      ts: 'number[]',
      zod: 'z.array(z.number())',
      control: 'text',
    },
    unique: {
      go: 'vectors.Vector',
      gorm: 'uniqueIndex',
      ts: 'number[]',
      zod: 'z.array(z.number()).optional()',
      control: 'text',
    },
    requiredUnique: {
      go: 'vectors.Vector',
      gorm: 'uniqueIndex',
      ts: 'number[]',
      zod: 'z.array(z.number())',
      control: 'text',
    },
    note: 'The dimension is part of the column type, so it cannot be changed later without a migration.',
  },
}
