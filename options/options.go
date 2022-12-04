package options

const TerraformMainFile string = "main.tf"
const PlanOutputFile = "plan.out"

var SoftLinkFiles = []string {"_accounts.tf", "_backend.tf", "_envs.tf", "_providers.tf"}

// terraform source code live under "provider" (aws/azure/gcp, etc) directories
var CloudProviders string = `
{
	"providers":[
        {
            "id": "aws",
            "description": "Amazon cloud service",
            "directory": "terraform/aws"
        },
        {
            "id": "azure",
            "description": "Microsoft Azure cloud service",
            "directory": "terraform/azure"
        },
        {
            "id": "gcp",
            "description":"Google cloud service",
            "directory": "terraform/gcp"
        }
	]
}
`

type Options struct {
	Help          bool     `short:"h" long:"help" description:"show help message"`
	Verbose       []bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
	Interactive   bool     `short:"i" long:"interactive" description:"Interactive Mode?"`
	TFaction      string   `short:"t" long:"tfaction" description:"Terraform Action - plan, apply, plan-destroy or apply-destroy" default:"plan" required:"false"`
	AutoApprove   bool     `short:"a" long:"approve" description:"Auto Approve Mode?" required:"false"`
	Workspace     string   `short:"w" long:"workspace" description:"Env/Workspace" required:"false"`
	Module        []string `short:"m" long:"modules" description:"List of Modules" default:"" required:"false"`
	Key           string   `short:"k" long:"key" description:"Token Key to checkout repo" default:"" required:"false"`
	Branch        string   `short:"b" long:"branch" description:"Merge Request Branch" default:"" required:"false"`
	GitOps        string   `short:"o" long:"gitops" description:"GitOps adhoc yaml" default:"" required:"false"`
	PreReq        bool     `short:"p" long:"prereq" description:"Build Pre-requisite base modules?" required:"false"`
	ComponentDir  string   `short:"c" long:"component" description:"Component directory name" default:"main"`
}

type Workspaces struct {
	Workspaces []Workspace `json:"workspaces"`
}

// terraform "workspace", bucket, switch iam roles, etc during build
type Workspace struct {
	ID                  string `json:"id"`
	AccountName         string `json:"account_name"`
	AccountID           string `json:"account_id"`
	WorkspaceIamRole    string `json:"workspace_iam_role"`
	WorkspaceR53IamRole string `json:"workspace_r53_iam_role"`
	WorkspaceSsmIamRole string `json:"workspace_ssm_iam_role"`
	Region              string `json:"region"`
	BucketRegion        string `json:"bucket_region"`
	Bucket              string `json:"bucket"`
	DynamoDb            string `json:"dynamodb"`
}

// TF Providers
type Providers struct {
	Providers []Provider `json:"providers"`
}

type Provider struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Directory   string `json:"directory"`
}



// terraform action
type TFActions struct {
	TFActions []TFAction `json:"tfactions"`
}

type TFAction struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

var ActionText string = `
{
    "tfactions":[
        {
            "id": "apply",
            "description": "Go ahead and perform terraform apply"
        },
        {
            "id": "apply-destroy",
            "description": "Go ahead and perform terraform destroy"
        },
        {
            "id": "plan",
            "description": "Dry-run to show me what's being changed"
        },
        {
            "id": "plan-destroy",
            "description": "Dry-run to show me what's being destroyed"
        }
    ]
}
`

// `BuildData` struct contains all info needed for terraform to run the build
type BuildData struct {
	Workspace    string
	Modules      []string
	AutoApprove  bool
	Interactive  bool
	TFaction     string
	Bucket       string
	BucketRegion string
	Dynamodb     string
	ProviderPath string
}