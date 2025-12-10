using System;
using System.IO;
using System.Threading.Tasks;
using Tachograph;

namespace Tachograph.CLI
{
    /// <summary>
    /// Command-line interface for the Tachograph SDK.
    /// </summary>
    class Program
    {
        static async Task<int> Main(string[] args)
        {
            if (args.Length == 0)
            {
                ShowHelp();
                return 0;
            }

            string command = args[0].ToLower();

            if (command == "parse")
            {
                return await ParseCommand(args);
            }
            else if (command == "help" || command == "--help" || command == "-h")
            {
                ShowHelp();
                return 0;
            }
            else
            {
                Console.Error.WriteLine($"Unknown command: {command}");
                Console.Error.WriteLine("Use 'tachograph help' for usage information.");
                return 1;
            }
        }

        static void ShowHelp()
        {
            Console.WriteLine("Tachograph CLI - Parse and analyze tachograph files");
            Console.WriteLine();
            Console.WriteLine("USAGE:");
            Console.WriteLine("  tachograph parse [options] <file> [<file>...]");
            Console.WriteLine();
            Console.WriteLine("COMMANDS:");
            Console.WriteLine("  parse              Parse .DDD files");
            Console.WriteLine("  help               Show this help message");
            Console.WriteLine();
            Console.WriteLine("OPTIONS:");
            Console.WriteLine("  --raw              Output raw intermediate format (skip semantic parsing)");
            Console.WriteLine("  --authenticate     Authenticate signatures and certificates");
            Console.WriteLine("  --strict           Error on unrecognized tags (default: true)");
            Console.WriteLine("  --no-strict        Don't error on unrecognized tags");
            Console.WriteLine("  --preserve-raw-data  Store raw bytes for round-trip fidelity (default: true)");
            Console.WriteLine();
        }

        static async Task<int> ParseCommand(string[] args)
        {
            bool raw = false;
            bool authenticate = false;
            bool strict = true;
            bool preserveRawData = true;
            var files = new System.Collections.Generic.List<string>();

            // Parse arguments
            for (int i = 1; i < args.Length; i++)
            {
                string arg = args[i];
                
                if (arg == "--raw")
                    raw = true;
                else if (arg == "--authenticate")
                    authenticate = true;
                else if (arg == "--strict")
                    strict = true;
                else if (arg == "--no-strict")
                    strict = false;
                else if (arg == "--preserve-raw-data")
                    preserveRawData = true;
                else if (arg.StartsWith("--"))
                {
                    Console.Error.WriteLine($"Unknown option: {arg}");
                    return 1;
                }
                else
                {
                    files.Add(arg);
                }
            }

            if (files.Count == 0)
            {
                Console.Error.WriteLine("Error: No files specified");
                Console.Error.WriteLine("Usage: tachograph parse [options] <file> [<file>...]");
                return 1;
            }

            // Process each file
            foreach (var filename in files)
            {
                try
                {
                    Console.WriteLine($"Processing: {filename}");

                    // Read file
                    if (!System.IO.File.Exists(filename))
                    {
                        Console.Error.WriteLine($"Error: File not found: {filename}");
                        return 1;
                    }

                    var data = await System.IO.File.ReadAllBytesAsync(filename);

                    // Step 1: Unmarshal to raw format
                    var unmarshalOpts = new UnmarshalOptions { Strict = strict };
                    var rawFile = unmarshalOpts.Unmarshal(data);
                    Console.WriteLine($"✓ Unmarshaled {filename}");

                    // Step 2: Optionally authenticate
                    if (authenticate)
                    {
                        var authOpts = new AuthenticateOptions { Mutate = true };
                        rawFile = await authOpts.AuthenticateAsync(rawFile);
                        Console.WriteLine($"✓ Authenticated {filename}");
                    }

                    // Step 3: Output raw or parse to semantic format
                    if (raw)
                    {
                        Console.WriteLine($"Raw file type: {rawFile.Type}");
                        // TODO: Output raw format (with or without authentication)
                    }
                    else
                    {
                        // Parse to semantic format
                        var parseOpts = new ParseOptions { PreserveRawData = preserveRawData };
                        var parsedFile = parseOpts.Parse(rawFile);
                        Console.WriteLine($"✓ Parsed {filename} (Type: {parsedFile.Type})");
                        
                        // TODO: Output parsed format (protobuf JSON)
                    }

                    Console.WriteLine();
                }
                catch (Exception ex)
                {
                    Console.Error.WriteLine($"Error processing {filename}: {ex.Message}");
                    Console.Error.WriteLine($"  {ex.GetType().Name}");
                    if (ex.InnerException != null)
                    {
                        Console.Error.WriteLine($"  Inner: {ex.InnerException.Message}");
                    }
                    return 1;
                }
            }

            Console.WriteLine($"Successfully processed {files.Count} file(s)");
            return 0;
        }
    }
}
