using System;

namespace Tachograph
{
    /// <summary>
    /// Configures the parsing process for converting raw tachograph files into semantic data structures.
    /// </summary>
    public class ParseOptions
    {
        /// <summary>
        /// Controls whether raw byte slices are stored in the RawData field of parsed protobuf messages.
        /// 
        /// If true (default), the raw byte slice for each parsed element will be stored in the RawData field,
        /// enabling perfect binary round-tripping via Marshal.
        /// 
        /// If false, RawData fields will be left empty, reducing memory usage but preventing exact binary reconstruction.
        /// </summary>
        public bool PreserveRawData { get; set; } = true;

        /// <summary>
        /// Performs semantic parsing on raw tachograph records, converting raw records into semantic data structures.
        /// If the raw file has been authenticated, the authentication results are propagated to the parsed messages.
        /// </summary>
        /// <param name="rawFile">The raw file to parse</param>
        /// <returns>A parsed File object</returns>
        /// <exception cref="ArgumentNullException">Thrown when rawFile is null</exception>
        /// <exception cref="NotSupportedException">Thrown when the file type is not supported</exception>
        public File Parse(RawFile rawFile)
        {
            if (rawFile == null)
                throw new ArgumentNullException(nameof(rawFile));

            switch (rawFile.Type)
            {
                case RawFileType.Card:
                    var cardType = InferCardType(rawFile.Card);
                    switch (cardType)
                    {
                        case CardType.DriverCard:
                            var driverCard = ParseDriverCardFile(rawFile.Card);
                            return new File
                            {
                                Type = FileType.DriverCard,
                                DriverCard = driverCard
                            };
                        default:
                            throw new NotSupportedException($"Unsupported card type: {cardType}");
                    }

                case RawFileType.VehicleUnit:
                    var vuFile = ParseVehicleUnitFile(rawFile.VehicleUnit);
                    return new File
                    {
                        Type = FileType.VehicleUnit,
                        VehicleUnit = vuFile
                    };

                default:
                    throw new NotSupportedException($"Unknown raw file type: {rawFile.Type}");
            }
        }

        private CardType InferCardType(RawCardFile card)
        {
            // TODO: Implement card type inference
            return CardType.DriverCard;
        }

        private DriverCardFile ParseDriverCardFile(RawCardFile card)
        {
            // TODO: Implement driver card file parsing
            throw new NotImplementedException("Driver card file parsing not yet implemented");
        }

        private VehicleUnitFile ParseVehicleUnitFile(RawVehicleUnitFile vu)
        {
            // TODO: Implement VU file parsing
            throw new NotImplementedException("Vehicle unit file parsing not yet implemented");
        }
    }
}
