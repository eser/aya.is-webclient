# Development Guidelines

This document contains essential rules and conventions for working with the AYA project codebase.

## Code Quality & Style

### Explicit Checks (CRITICAL)
- **NEVER use implicit/truthy/falsy checks** except for boolean values
- ✅ **Correct**: `if (value === null)`, `if (string === "")`, `if (array.length === 0)`
- ❌ **Wrong**: `if (!value)`, `if (!string)`, `if (!array.length)`
- ✅ **Exception**: Boolean checks are fine: `if (!response.ok)`, `if (isValid)`
- ✅ **Ternary**: Use explicit checks: `value !== null ? value : defaultValue`

### Backend Function Organization
- **Always create dedicated backend functions** in `/webclient/src/shared/modules/backend/`
- Follow existing patterns like `get-current-session.ts`
- **Never inline fetch logic** in components
- Export functions from index files for clean imports
- Use proper TypeScript interfaces for request/response types

### UI Component Usage & Styling (CRITICAL METHODOLOGY)

#### CSS Modules with @apply (STRONGLY PREFERRED METHOD)
- **Primary styling approach**: CSS Modules with Tailwind `@apply` directive
- **Each component** should have a co-located `*.module.css` file
- **Semantic class names** that describe element's role (e.g., `.product-card`, `.card-title`)
- **Direct Tailwind classes** should be "used sparingly and are generally discouraged"

```css
/* Example: product-card.module.css */
.product-card {
  @apply border rounded-lg p-4 shadow-md flex flex-col gap-2;

  & .title {
    @apply text-xl font-bold mb-2;
  }

  & .description {
    @apply text-sm text-gray-600;
  }
}
```

```tsx
// Example: ProductCard.tsx
import styles from './product-card.module.css';

type ProductCardProps = {
  product: { name: string; description: string; };
};

function ProductCard(props: ProductCardProps) {
  return (
    <div className={styles['product-card']}>
      <h3 className={styles.title}>{props.product.name}</h3>
      <p className={styles.description}>{props.product.description}</p>
    </div>
  );
}
```

#### When Direct Tailwind is Acceptable (RARE)
- **Very simple, non-reusable micro-adjustments** on generic HTML elements
- **Global layout containers** in top-level pages (e.g., `layout.tsx`)
- **Use `cn` utility** when combining: `className={cn(styles.button, "mt-2")}`

#### Component Standards
- **Always use Shadcn UI components** when available
- **Never use inline styles** - violates maintainability
- **Single props object**: `function Component(props: ComponentProps)` not destructuring
- Add `type` attribute to all buttons
- Use proper semantic HTML with accessibility

### Translation System & Internationalization

#### Translation Key Format
- Use English text as translation keys (not snake_case)
- ✅ **Correct**: `t("Auth", "Login with GitHub")`
- ❌ **Wrong**: `t("Auth.login_with_github")`
- Add translations to ALL language files in `/webclient/src/messages/`
- Use fallbacks for missing translations: `t("Section", "Key") || "Fallback"`

#### React Server Components vs Client Components
- **RSC**: Use `getTranslations` async function
- **Client Components**: Use `useTranslations` hook with `TranslationsProvider`

```tsx
// ✅ React Server Component
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";

async function MyServerComponent() {
  const { t, locale } = await getTranslations();
  return <h1>{t("Home", "AYA the Open Source Network")}</h1>;
}

// ✅ Client Component
"use client";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";

function MyClientComponent() {
  const { t, locale } = useTranslations();
  return <button>{t("Layout", "Change theme")}</button>;
}
```

#### Message File Management
- **Location**: `src/messages/[locale-code].json`
- **Structure**: Nested JSON with consistent keys across all files
- **Reference**: `en.json` is the master reference for all keys
- **Adding new strings**:
  1. Add to `en.json` first with English value
  2. Copy keys to all other language files
  3. Provide translations or leave English as placeholder

## Project Architecture

### Centralized Backend Object Pattern (CRITICAL)
- **Primary backend communication** goes through centralized module: `src/shared/modules/backend/backend.ts`
- **Import the backend object**: `import { backend } from "@/shared/modules/backend/backend.ts"`
- **Access methods via backend object**: `backend.getProfile("en", profileId)`
- **Single source of truth** for all backend functionality
- **Type imports only** from submodules: `import type { MyType } from '@/shared/modules/backend/profiles/types'`
- **Benefits**: Clear contract, simplified interaction, abstracted internal structure

```typescript
// ✅ Correct backend usage
import { backend } from "@/shared/modules/backend/backend.ts";

async function UserProfile(props: { profileId: string }) {
  const profile = await backend.getProfile("en", props.profileId);
  return <div>{profile.title}</div>;
}
```

### Frontend Webclient Structure
```
webclient/src/
├── app/                          # Next.js pages (app router)
├── shared/
│   ├── modules/backend/          # Centralized backend object + submodules
│   ├── components/ui/            # Shadcn components
│   ├── modules/auth/             # Authentication logic
│   └── modules/i18n/             # Translation system
└── messages/                     # Translation files (13 languages)
```

### Backend Services Structure (Hexagonal Architecture)
```
services/pkg/
├── api/
│   ├── business/                 # Pure business logic (ports/interfaces)
│   └── adapters/                 # External implementations (HTTP, storage)
│       ├── http/                 # HTTP routes and handlers
│       ├── storage/              # Database implementations
│       └── auth_providers/       # Authentication adapters
├── ajan/                         # Framework utilities
└── lib/                          # Shared libraries
```

### Development Commands
```bash
# Frontend (in /webclient directory)
deno task lint        # Run Deno linter
deno task lint:fix    # Fix linting issues automatically
deno task format      # Check formatting
deno task format:fix  # Fix formatting issues
deno task typecheck   # Run TypeScript type checking (if available)

# Backend (in /services directory)
make restart         # Restart services after changes
make check           # Run static analysis tools
make lint            # Run linting
make fix             # Fix formatting and linting issues
make test            # Run tests
```

## Hexagonal Architecture Principles (CRITICAL FOR SERVICES)

### Core Principles
- **Business logic is sacred** - keep it clean and dependency-free
- **All external interactions** must go through well-defined interfaces (ports)
- **Business logic packages** define interfaces and types first
- **Adapters implement** these business interfaces
- **No circular dependencies** allowed
- **Dependency injection** used to wire everything together

### Dependency Rules
- **Business logic** can only depend on other business logic
- **Adapters** can only depend on business interfaces they implement
- **Entry points** (cmd/) wire everything together
- **External dependencies** only in adapters, never in business logic

### File Structure Conventions
- Use **snake_case** for file names in services
- **Business errors** defined in business layer
- **Wrap external errors** before returning them
- **Test files** co-located with `_test` package suffix
- **Configuration** injected via appcontext

### Example Interface Pattern
```go
// In business layer: pkg/api/business/profiles/service.go
type Repository interface {
    GetProfile(ctx context.Context, id string) (*Profile, error)
    CreateProfile(ctx context.Context, profile *Profile) error
}

type Service struct {
    repo Repository // Interface, not concrete implementation
}

// In adapter: pkg/api/adapters/storage/repository_profiles.go
type ProfileRepository struct {
    db *sql.DB
}

func (r *ProfileRepository) GetProfile(ctx context.Context, id string) (*Profile, error) {
    // Implementation details
}
```

## Authentication & Security
- JWT tokens stored in localStorage (TODO: migrate to httpOnly cookies)
- Use authentication middleware for protected routes
- Follow existing session management patterns
- Never commit secrets or expose sensitive data

## Database & Migrations
- Use SQLC for type-safe database queries
- Create migrations in `/services/etc/data/default/migrations/`
- Follow existing naming conventions for tables and fields
- Use appropriate PostgreSQL data types

## Development Workflow & Quality Gates (MANDATORY)

### Testing & Linting (ALWAYS REQUIRED)
- **ALWAYS run tests and linters when making changes** - no exceptions
- **Fix linter errors immediately** - never leave them unfixed
- **Never commit** without passing all quality gates

### Testing Requirements (CRITICAL)
- **Business logic MUST have unit tests** with mock adapters
- **Adapters should have integration tests**
- **Write testable code** using dependency injection
- **Use interfaces for mocking** in tests (hexagonal architecture)
- **Test files** co-located with `_test` package suffix

#### Example Unit Test with Mock
```go
// profiles_test.go
type mockRepository struct{}

func (m *mockRepository) GetProfile(ctx context.Context, id string) (*Profile, error) {
    // Mock implementation
    return &Profile{ID: id, Title: "Test"}, nil
}

func TestProfileService_GetProfile(t *testing.T) {
    service := &Service{repo: &mockRepository{}}
    profile, err := service.GetProfile(context.Background(), "test-id")
    // Assert results
}
```

### Quality Gates Checklist
- [ ] All linters pass (`deno task lint`)
- [ ] All formatters pass (`deno task format`)
- [ ] TypeScript compilation successful
- [ ] All tests pass (unit + integration)
- [ ] No TODO comments or placeholders left
- [ ] Endpoints tested with curl/manual testing

### Git Workflow
- **Never commit without explicit request**
- Create descriptive commit messages
- Follow existing commit message patterns
- Use conventional commit format when possible

### Error Handling
- Use static error variables instead of dynamic strings
- Provide meaningful error messages to users
- Log errors appropriately with context
- Handle edge cases explicitly

## Specific Conventions

### Profile System
- Users can have only ONE individual profile
- Support organization and product profile types
- Use slug-based routing for profiles
- Implement real-time slug validation

### API Routes
- Use locale-based routing: `/{locale}/endpoint`
- Wrap responses in `{data: ..., error: null}` format for frontend compatibility
- Use proper HTTP status codes
- Implement authentication middleware where needed

### Form Validation
- Use Zod for form validation with React Hook Form
- Use `zodResolver` for integration
- Provide real-time validation feedback
- Handle server-side validation errors

### Internationalization
- Support 13 languages: ar, de, en, es, fr, it, ja, ko, nl, pt-PT, ru, tr, zh-CN
- Use `useTranslations` hook for client components
- Default to English (`en`) for backend API calls
- Detect locale from context, not URL parsing

## React & Component Patterns

### Component Props Pattern (REQUIRED)
- **Accept single props object** instead of destructuring in function signature
- **Access via `props.propertyName`** for better readability and refactoring
- **Define explicit prop types** as `ComponentNameProps`

```tsx
// ✅ Correct - Single props object
type UserProfileProps = {
  userId: string;
  showActions: boolean;
};

function UserProfile(props: UserProfileProps) {
  return (
    <div>
      <h1>User: {props.userId}</h1>
      {props.showActions && <button>Edit</button>}
    </div>
  );
}

// ❌ Wrong - Destructured props
function UserProfile({ userId, showActions }: UserProfileProps) {
  // Avoid this pattern
}
```

### Export Patterns
- **Prefer named exports** for React components: `export function MyComponent() {}`
- **Default exports only** for Next.js required files (page.tsx, layout.tsx)
- **Avoid** `export { MyComponent as default }` pattern

### Early Returns & Readability
- **Use early returns** whenever possible for better readability
- **Functional and declarative** patterns; avoid classes
- **Iteration and modularization** over code duplication

## Common Patterns

### Backend Function Template (via Central Object)
```typescript
// ❌ Wrong - Direct function import
import { getFunctionName } from "@/shared/modules/backend/some-module.ts";

// ✅ Correct - Via central backend object
import { backend } from "@/shared/modules/backend/backend.ts";
import type { ResponseType } from "@/shared/modules/backend/types.ts";

async function MyComponent(props: { locale: string; param: string }) {
  const response = await backend.getFunctionName(props.locale, props.param);
  return <div>{response?.title}</div>;
}
```

### Component Styling Pattern (CSS Modules)
```tsx
// ProductCard.tsx
import styles from './product-card.module.css';
import { Card, CardHeader, CardTitle, CardContent } from "@/shared/components/ui/card.tsx";

type ProductCardProps = {
  product: { title: string; description: string; };
};

function ProductCard(props: ProductCardProps) {
  return (
    <div className={styles['product-card']}>
      <Card>
        <CardHeader>
          <CardTitle className={styles.title}>{props.product.title}</CardTitle>
        </CardHeader>
        <CardContent className={styles.content}>
          <p>{props.product.description}</p>
        </CardContent>
      </Card>
    </div>
  );
}
```

## React v19 & Next.js Specific Rules

### React v19 Compiler
- **Write compiler-friendly code** - favor idiomatic React patterns
- **Minimize manual memoization** (`useMemo`, `useCallback`) unless profiling shows benefit
- **React compiler handles optimization** automatically
- **Avoid patterns** that might confuse the compiler

### Next.js Configuration
- **PPR enabled** (`experimental.ppr: true`) - understand rendering/caching implications
- **React Strict Mode** enabled for better development experience
- **Trailing slash** set to `false`
- **App Router** - use server components as default

### Server Components & "server-only"
- **Embrace React Server Components (RSC)** as the default
- **Use "server-only" package** for strict server-side modules
- **Async Server Components** for data fetching operations
- **Minimize "use client"** - favor RSC when possible
- **Wrap client components in Suspense** with fallback

### Metadata API
- **Use Next.js Metadata API** for `<head>` elements
- **Export metadata** in layouts and pages: `export const metadata = ...`
- **Export viewport** when needed: `export const viewport = ...`

### Middleware
- **Custom domains** and **internationalization** handled via middleware
- **Configuration** defined in `src/middleware.ts`
- **Runtime and matcher** properly configured

## Tooling & Environment

### Deno Configuration
- **TypeScript Strict Mode** enabled via `deno.json` and `tsconfig.json`
- **Linting**: Uses Deno's `recommended` rules
- **Formatting**: Line width 120 characters, 2 spaces indentation
- **Path Aliases**: Use `@/` for imports from `src` directory

### Project Configuration
- **Pre-commit hooks** ensure quality gates pass
- **Environment variables** for runtime configuration
- **`.env.local`** and **`config.local.json`** for development overrides
- **Keep secrets out** of codebase

## Database & Code Generation

### SQLC Patterns
- **All SQL queries** in `etc/data/{datasource_name}/queries/`
- **Database code generated** using SQLC
- **Clean separation** between generated and hand-written code
- **Migrations** in `etc/data/default/migrations/` (forward-only)

### Schema Management
- **Document schema changes**
- **Use appropriate PostgreSQL data types**
- **Follow existing naming conventions** for tables and fields

## Debugging & Development
- Use browser dev tools for frontend debugging
- Check network tab for API call issues
- Use `console.error` for error logging (remove before production)
- Test with different locales and user states
- **Manual endpoint testing** with curl before committing

## Performance
- Use React Server Components where appropriate
- Implement proper loading states
- Use TypeScript for better development experience
- Minimize bundle size with proper imports
- **Optimize Web Vitals** (LCP, CLS, FID)

## Remember
- **EXPLICIT IS BETTER THAN IMPLICIT** - especially for null/undefined checks
- **Business logic is sacred** - keep it clean and dependency-free
- **Always run linters and fix errors immediately**
- **Follow hexagonal architecture** for backend services
- **Use CSS Modules with @apply** for styling
- **Single props object** for components
- **Central backend object** for API calls
- Follow existing codebase patterns religiously
- Test thoroughly before committing
